type AppleParams struct {
	IosVersion    string `json:"ios_version" desc:"ios版本"`
	OrderId       string `json:"order_id" valid:"Required" desc:"订单号"`
	TransactionId string `json:"transaction_id" valid:"Required" desc:"支付流水"`
	Receipt       string `json:"receipt" valid:"Required" desc:""`
}

// @Title 苹果支付回调接口
// @Description 苹果支付回调接口
// @Param   body       body   pay.AppleParams  true
// @Success 0 {string} 状态码
// @router /apple-callback [post]
func (b *PayController) AppleCallBack() {
	var p AppleParams
	if !b.DecodeParams(&p) || !b.ValidParams(&p) {
		return
	}
	order, err := pay.PayOrderByPayId(p.OrderId)
	if err != nil || order.Uid != b.Uid {
		b.ResponseErrorMsg("-1", "订单不存在")
		return
	}
	if order.TransactionId != "" && order.TransactionId != p.TransactionId {
		log.Errorf("订单信息不匹配:%#v--%#v", order, p)
		b.ResponseErrorMsg("-1", "订单信息不匹配")
		return
	}
	requestData := map[string]interface{}{
		"receipt-data": p.Receipt,
	}
	body, err := json.Marshal(requestData)
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	var productId uint64     //苹果产品id
	var transactionId string //苹果订单号
	//提交到苹果正式服务器验证
	proUrl := fmt.Sprintf("%v", AppleProVerifyUrl)
	rsp, err := fetch.Cmd(fetch.Request{
		Method: "POST",
		URL:    proUrl,
		Body:   body,
		Header: http.Header{"Content-Type": {"application/json"}},
	})
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	log.Debugf(fmt.Sprintf("苹果正式服务器返回值:%#v\n", string(rsp)))
	r, err := apple_pay.Verify(rsp)
	if err != nil {
		log.Errorf(fmt.Sprintf("解析出错:%#v\n", err))
		b.ResponseErrorMsg("-1", "系统异常")
		return
	}
	//提交到苹果测试服务器验证
	if r.Status == 21007 {
		testUrl := fmt.Sprintf("%v", AppleTestVerifyUrl)
		rsp, err := fetch.Cmd(fetch.Request{
			Method: "POST",
			URL:    testUrl,
			Body:   body,
			Header: http.Header{"Content-Type": {"application/json"}},
		})
		if err != nil {
			log.Errorf(err.Error())
			return
		}
		log.Debugf(fmt.Sprintf("苹果测试服务器返回值:%#v\n", string(rsp)))
		r, err := apple_pay.Verify(rsp)
		if err != nil {
			log.Errorf(fmt.Sprintf("解析出错:%#v\n", err))
			b.ResponseErrorMsg("-1", "系统异常1")
			return
		}
		if r.Status != 0 {
			log.Errorf(fmt.Sprintf("苹果测试服务器返回错误,错误码:%v\n", r.Status))
			b.ResponseErrorMsg("-1", "系统异常2")
			return
		}
		if len(r.Receipt.InApp) < 1 {
			log.Errorf(fmt.Sprintf("苹果正式服务器返回数据InApp出错"))
			b.ResponseErrorMsg("-1", "系统异常3")
			return
		}

		productId, _ = strconv.ParseUint(r.Receipt.InApp[0].ProductId, 10, 64)
		transactionId = r.Receipt.InApp[0].TransactionId

	} else if r.Status != 0 {
		log.Errorf(fmt.Sprintf("苹果正式服务器返回错误,第一错误码:%v\n", r.Status))
		b.ResponseErrorMsg("-1", "系统异常4")
		return
	} else {
		productId, _ = strconv.ParseUint(r.Receipt.InApp[0].ProductId, 10, 64)
		transactionId = r.Receipt.InApp[0].TransactionId
	}
	if productId != order.ProductId {
		b.ResponseErrorMsg("-1", "订单信息错误")
	}
	if pay.CheckAppleTransactionIdIsPay(transactionId) > 0 {
		b.ResponseErrorMsg("-1", "重复伪造订单支付")
		return
	}
	//----------------------下面是支付成功的情况----------------------
	pay.PaySuccess(p.OrderId, transactionId, order.Price, "", "", false)
	b.ResponseSuccess()
}