package sdk

import (
	"encoding/json"
	"fmt"
	"log"
)

// Client 客户端
type Client struct {
	BrokerID string // 代征主体ID
	DealerID string // 商户ID
	Gateway  string // 路由
	Appkey   string // 商户appkey
	Des3Key  string // 商户des3key
}

// New 新建客户端
func New(brokerID, dealerID, gateway, appkey, des3key string) *Client {
	return &Client{
		BrokerID: brokerID,
		DealerID: dealerID,
		Gateway:  gateway,
		Appkey:   appkey,
		Des3Key:  des3key,
	}
}

// CreateBankOrder 创建银行卡订单
func (c *Client) CreateBankOrder(param *BankOrderParam) (ref string, err error) {
	fn := "CreateBankOrder"
	order := map[string]string{
		"order_id":   param.OrderID,
		"broker_id":  c.BrokerID,
		"dealer_id":  c.DealerID,
		"real_name":  param.RealName,
		"id_card":    param.IDCard,
		"card_no":    param.CardNo,
		"phone_no":   param.PhoneNo,
		"pay":        param.Pay,
		"pay_remark": param.PayRemark,
		"notify_url": param.NotifyURL,
	}

	params, err := BuildParams(order, c.Appkey, c.Des3Key)
	if err != nil {
		log.Printf("@[%s] BuildParams failed, err=%v, order=%+v", fn, err, order)
		return
	}

	headers := BuildHeader(c.DealerID)

	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Post(c.Gateway+BankOrderURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Post failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, params=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res CreateOrderResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, result=%s, params=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == SuccessCode {
		ref = res.Data.Ref
		return
	}

	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// CreateAliOrder 创建支付宝订单
func (c *Client) CreateAliOrder(param *AliOrderParam) (ref string, err error) {
	fn := "CreateAliOrder"
	order := map[string]string{
		"order_id":   param.OrderID,
		"broker_id":  c.BrokerID,
		"dealer_id":  c.DealerID,
		"real_name":  param.RealName,
		"id_card":    param.IDCard,
		"card_no":    param.CardNo,
		"pay":        param.Pay,
		"pay_remark": param.PayRemark,
		"check_name": param.CheckName,
		"notify_url": param.NotifyURL,
	}

	params, err := BuildParams(order, c.Appkey, c.Des3Key)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, order=%+v", fn, err, order)
		return
	}

	headers := BuildHeader(c.DealerID)

	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Post(c.Gateway+AliOrderURL, params, headers)
	if err != nil {
		log.Printf("@[%s] util.Post failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res CreateOrderResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == SuccessCode {
		ref = res.Data.Ref
		return
	}

	err = fmt.Errorf("Err:%s", res.Message)
	return
}

//CereateWxOrder 创建微信订单
func (c *Client) CereateWxOrder(param *WxOrderParam) (ref string, err error) {
	fn := "CereateWxOrder"
	order := map[string]string{
		"order_id":   param.OrderID,
		"broker_id":  c.BrokerID,
		"dealer_id":  c.DealerID,
		"real_name":  param.RealName,
		"id_card":    param.IDCard,
		"openid":     param.OpenID,
		"pay":        param.Pay,
		"notes":      param.Notes,
		"pay_remark": param.PayRemark,
		"notify_url": param.NotifyURL,
		"wx_app_id":  param.WxAppID,
		"wxpay_mode": param.WxPayMode,
	}

	params, err := BuildParams(order, c.Appkey, c.Des3Key)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, order=%+v", fn, err, order)
		return
	}

	headers := BuildHeader(c.DealerID)

	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Post(c.Gateway+WxOrderURL, params, headers)
	if err != nil {
		log.Printf("@[%s] util.Post failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res CreateOrderResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == SuccessCode {
		ref = res.Data.Ref
		return
	}

	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// OrderCallBack 订单回调
func (c *Client) OrderCallBack(data, mess, timestamp, sign string) (order OrderDetailInfo, err error) {
	fn := "OrderCallBack"
	log.Printf("@[%s] begin, data=%s, mess=%s, timestamp=%s, sign=%s", fn, data, mess, timestamp, sign)
	originSign := Sign(data, mess, timestamp, c.Appkey)
	if originSign != sign {
		log.Printf("@[%s] sign mismatch, data=%s, mess=%s, timestamp=%s, originSign=%s, sign=%s", fn, data, mess, timestamp, originSign, sign)
	}

	originData, err := Decrypt(data, c.Des3Key)
	if err != nil {
		log.Printf("@[%s] Decrypt failed, err=%v, data=%s", fn, err, data)
		return
	}

	var res OrderCallBackResponse

	err = json.Unmarshal(originData, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, data=%s", fn, err, string(originData))
		return
	}
	log.Printf("@[%s] end, res=%+v", fn, res)
	order = res.Data
	return
}

// QueryOrder 查询订单信息
func (c *Client) QueryOrder(orderID, channel, dataType string) (dest OrderInfo, err error) {
	fn := "QueryOrder"
	data := map[string]string{
		"order_id":  orderID,
		"channel":   channel,
		"data_type": dataType,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Get(c.Gateway+QueryOrderURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Get failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res QueryOrderResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == SuccessCode {
		dest = res.Data
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// QueryAccountBalance 查询账户余额
func (c *Client) QueryAccountBalance() (accounts []AccountBalance, err error) {
	fn := "QueryAccountBalance"
	data := map[string]string{
		"dealer_id": c.DealerID,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Get(c.Gateway+QueryAccountURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Get failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res QueryAccountBalanceResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == SuccessCode {
		accounts = res.Data.DealerInfos
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// QueryReceiptFile 查询电子回单
func (c *Client) QueryReceiptFile(orderID, ref string) (file OrderReceiptFile, err error) {
	fn := "QueryOrderReceiptFile"
	data := map[string]string{
		"order_id": orderID,
		"ref":      ref,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Get(c.Gateway+QueryReceiptFileURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Get failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res QueryReceiptFileResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == SuccessCode {
		file = res.Data
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// CancelOrder 取消订单
func (c *Client) CancelOrder(orderID, ref, channel string) (ok bool, err error) {
	fn := "CancelOrder"
	data := map[string]string{
		"dealer_id": c.DealerID,
		"order_id":  orderID,
		"ref":       ref,
		"channel":   channel,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Post(c.Gateway+CancelOrderURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Post failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res BaseCheckResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == SuccessCode {
		ok = res.Data.Ok
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// DownloadOrderFile 下载日订单文件
func (c *Client) DownloadOrderFile(orderDate string) (url string, err error) {
	fn := "DownloadOrderFile"
	data := map[string]string{
		"order_date": orderDate,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Get(c.Gateway+DownloadOrderURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Get failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res DownloadOrderFileResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == SuccessCode {
		url = res.Data.OrderDownloadURL
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// DownloadBillFile 下载日流水文件
func (c *Client) DownloadBillFile(billDate string) (url string, err error) {
	fn := "DownloadBillFile"
	data := map[string]string{
		"bill_date": billDate,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Get(c.Gateway+DownloadBillURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Get failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res DownloadBillFileResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == SuccessCode {
		url = res.Data.BillDownloadURL
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// QueryRechargeRecord 查询充值记录
func (c *Client) QueryRechargeRecord(beginAt, endAt string) (records []RechargeRecord, err error) {
	fn := "QueryRechargeRecord"
	data := map[string]string{
		"begin_at": beginAt,
		"end_at":   endAt,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Get(c.Gateway+QueryRechargeURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Get failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res QueryRechargeRecordResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == SuccessCode {
		records = res.Data
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// UploadUserInfo 上传免验证用户名单信息
func (c *Client) UploadUserInfo(param *UserInfoParam) (ok bool, err error) {
	fn := "UploadUserInfo"
	data := map[string]interface{}{
		"broker_id":     c.BrokerID,
		"dealer_id":     c.DealerID,
		"ref":           param.Ref,
		"id_card":       param.IDCard,
		"real_name":     param.RealName,
		"card_type":     param.CardType,
		"country":       param.Country,
		"birthday":      param.Birthday,
		"gender":        param.Gender,
		"user_images":   param.UserImages,
		"comment_apply": param.CommentApply,
		"notify_url":    param.NotifyURL,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Post(c.Gateway+UploadUserURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Post failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res BaseCheckResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == SuccessCode {
		ok = res.Data.Ok
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// UserInfoCallback 免验证用户信息上传回调信息
func (c *Client) UserInfoCallback(data, mess, timestamp, sign string) (user UserCallBackInfo, err error) {
	fn := "UserInfoCallback"
	log.Printf("@[%s] begin, data=%s, mess=%s, timestamp=%s, sign=%s", fn, data, mess, timestamp, sign)
	originSign := Sign(data, mess, timestamp, c.Appkey)
	if originSign != sign {
		log.Printf("@[%s] sign mismatch, data=%s, mess=%s, timestamp=%s, originSign=%s, sign=%s", fn, data, mess, timestamp, originSign, sign)
	}

	originData, err := Decrypt(data, c.Des3Key)
	if err != nil {
		log.Printf("@[%s] Decrypt failed, err=%v, data=%s", fn, err, data)
		return
	}

	var res UserCallBackInfo

	err = json.Unmarshal(originData, &user)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, data=%s", fn, err, string(originData))
		return
	}
	log.Printf("@[%s] end, res=%+v", fn, res)
	return
}

// CheckUserExist 校验免验证用户是否存在
func (c *Client) CheckUserExist(idCard, realName string) (ok bool, err error) {
	fn := "CheckUserExist"
	data := map[string]string{
		"id_card":   idCard,
		"real_name": realName,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Post(c.Gateway+CheckExistUserURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Post failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res BaseCheckResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == SuccessCode {
		ok = res.Data.Ok
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// QueryInvoice 查询发票信息
func (c *Client) QueryInvoice(year int) (invoice InvoiceInfo, err error) {
	fn := "QueryInvoice"
	data := map[string]interface{}{
		"broker_id": c.BrokerID,
		"dealer_id": c.DealerID,
		"year":      year,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Get(c.Gateway+QueryInvoiceURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Get failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res QueryInvoiceResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == SuccessCode {
		invoice = res.Data
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// ElementVerifyRequest 银行卡四要素请求鉴权
func (c *Client) ElementVerifyRequest(idCard, realName, cardNo, mobile string) (ref string, err error) {
	fn := "ElementVerifyRequest"
	data := map[string]string{
		"id_card":   idCard,
		"real_name": realName,
		"card_no":   cardNo,
		"mobile":    mobile,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Post(c.Gateway+Element4RequestURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Post failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res ElementVerifyResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == SuccessCode {
		ref = res.Data.Ref
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// ElementVerifyConfirm 银行卡四要素确认鉴权
func (c *Client) ElementVerifyConfirm(idCard, realName, cardNo, mobile, ref, captcha string) (ok bool, err error) {
	fn := "ElementVerifyConfirm"
	data := map[string]string{
		"id_card":   idCard,
		"real_name": realName,
		"card_no":   cardNo,
		"mobile":    mobile,
		"ref":       ref,
		"captcha":   captcha,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Post(c.Gateway+Element4ConfirmURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Post failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res BaseResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == SuccessCode {
		ok = true
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// Element4Check 银行卡四要素鉴权
func (c *Client) Element4Check(idCard, realName, cardNo, mobile string) (ok bool, err error) {
	fn := "Element4Check"
	data := map[string]string{
		"id_card":   idCard,
		"real_name": realName,
		"card_no":   cardNo,
		"mobile":    mobile,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Post(c.Gateway+Element4URL, params, headers)
	if err != nil {
		log.Printf("@[%s] Post failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res BaseResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == SuccessCode {
		ok = true
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// Element3Check 银行卡三要素鉴权
func (c *Client) Element3Check(idCard, realName, cardNo string) (ok bool, err error) {
	fn := "Element3Check"
	data := map[string]string{
		"id_card":   idCard,
		"real_name": realName,
		"card_no":   cardNo,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Post(c.Gateway+Element3URL, params, headers)
	if err != nil {
		log.Printf("@[%s] Post failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res BaseResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == SuccessCode {
		ok = true
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// IDCheck 实名制二要素鉴权接口
func (c *Client) IDCheck(idCard, realName string) (ok bool, err error) {
	fn := "IDCheck"
	data := map[string]string{
		"id_card":   idCard,
		"real_name": realName,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Post(c.Gateway+IDCheckURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Post failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res BaseResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == SuccessCode {
		ok = true
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}

// QueryBankCardInfo 查询银行卡信息
func (c *Client) QueryBankCardInfo(cardNo, bankName string) (cardInfo BankCardInfo, err error) {
	fn := "QueryBankCardInfo"
	data := map[string]string{
		"card_no":   cardNo,
		"bank_name": bankName,
	}

	params, err := BuildParams(data, c.Appkey, c.Des3Key)
	if err != nil {
		log.Printf("@[%s] BuildParam failed, err=%v, data=%+v", fn, err, data)
		return
	}

	headers := BuildHeader(c.DealerID)
	log.Printf("@[%s] begin, params=%+v, header=%v", fn, params, headers)

	result, err := Get(c.Gateway+BankCardInfoURL, params, headers)
	if err != nil {
		log.Printf("@[%s] Get failed, err=%v, params=%+v, headers=%+v", fn, err, params, headers)
		return
	}
	log.Printf("@[%s] end, data=%+v, header=%v, result=%s", fn, params, headers, string(result))

	var res QueryBankCardResponse
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Printf("@[%s] json.Unmarshal failed, err=%v, params=%s, data=%+v, header=%+v", fn, err, string(result), params, headers)
		return
	}

	// 请求成功
	if res.Code == SuccessCode {
		cardInfo = res.Data
		return
	}
	err = fmt.Errorf("Err:%s", res.Message)
	return
}
