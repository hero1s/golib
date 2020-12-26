package gopay

import (
	"git.moumentei.com/plat_go/golib/helpers/ip"
	"git.moumentei.com/plat_go/golib/log"
	"github.com/iGoogle-ink/gopay"
	"strconv"
	"time"
)

type WeChatPayClient struct {
	PayClient      *gopay.WeChatClient
	JsApiPayClient *gopay.WeChatClient
	WechatPay      WeChatPayParam `json:"wechat_pay" desc:"微信支付参数"`
	WechatJsPay    WeChatPayParam `json:"wechat_js_pay" desc:"微信JSAPI支付参数"`
}

func InitWechatPay(isProd bool, pay, payJs WeChatPayParam) *WeChatPayClient {
	client := new(WeChatPayClient)

	client.PayClient = gopay.NewWeChatClient(pay.WeChatAppId, pay.WeChatMchId, pay.WeChatKey, isProd)
	client.PayClient.SetCountry(gopay.China)
	client.WechatPay = pay

	client.JsApiPayClient = gopay.NewWeChatClient(payJs.WeChatAppId, payJs.WeChatMchId, payJs.WeChatKey, isProd)
	client.JsApiPayClient.SetCountry(gopay.China)
	client.WechatJsPay = payJs

	return client
}
func (c *WeChatPayClient) getPayClient(tradeType string) *gopay.WeChatClient {
	if tradeType == gopay.TradeType_JsApi {
		return c.JsApiPayClient
	}
	return c.PayClient
}
func (c *WeChatPayClient) getPayParam(tradeType string) WeChatPayParam {
	if tradeType == gopay.TradeType_JsApi {
		return c.WechatJsPay
	}
	return c.WechatPay
}

//微信预下单
func (c *WeChatPayClient) UnifiedOrder(moneyFee int64, describe, orderId, tradeType, deviceInfo, openid string) (map[string]string, error) {
	//初始化参数Map
	body := make(gopay.BodyMap)
	body.Set("nonce_str", gopay.GetRandomString(32))
	body.Set("body", describe)
	body.Set("out_trade_no", orderId)
	body.Set("total_fee", moneyFee) //单位分
	body.Set("spbill_create_ip", ip.LocalIP())
	body.Set("notify_url", c.getPayParam(tradeType).WeChatCallbackUrl)
	body.Set("trade_type", tradeType)
	body.Set("device_info", deviceInfo)
	body.Set("sign_type", gopay.SignType_MD5)

	//请求支付下单，成功后得到结果
	var cp = make(map[string]string)
	if tradeType == gopay.TradeType_JsApi {
		body.Set("openid", openid)
	}
	wxRsp, err := c.getPayClient(tradeType).UnifiedOrder(body)
	if err != nil {
		log.Errorf("微信预下单:%#v  \n支付失败Error:%v", body, err.Error())
		return cp, err
	} else {
		log.Debugf("微信预下单wxRsp:%#v", *wxRsp)
	}

	cp["appid"] = wxRsp.Appid
	cp["partnerid"] = wxRsp.MchId
	cp["prepayid"] = wxRsp.PrepayId
	cp["package"] = "Sign=WXPay"
	cp["noncestr"] = wxRsp.NonceStr
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	if tradeType == gopay.TradeType_App {
		sign := gopay.GetAppPaySign(wxRsp.Appid, wxRsp.MchId, wxRsp.NonceStr, wxRsp.PrepayId, gopay.SignType_MD5, timeStamp, c.getPayParam(tradeType).WeChatKey)
		cp["paySign"] = sign
	} else if tradeType == gopay.TradeType_JsApi {
		pac := "prepay_id=" + wxRsp.PrepayId
		sign := gopay.GetMiniPaySign(wxRsp.Appid, wxRsp.NonceStr, pac, gopay.SignType_MD5, timeStamp, c.getPayParam(tradeType).WeChatKey)
		cp["paySign"] = sign
	} else if tradeType == gopay.TradeType_H5 {
		pac := "prepay_id=" + wxRsp.PrepayId
		sign := gopay.GetH5PaySign(wxRsp.Appid, wxRsp.NonceStr, pac, gopay.SignType_MD5, timeStamp, c.getPayParam(tradeType).WeChatKey)
		cp["paySign"] = sign
		cp["mweb_url"] = wxRsp.MwebUrl
	} else if tradeType == gopay.TradeType_Native {
		cp["code_url"] = wxRsp.CodeUrl
	}

	cp["timestamp"] = timeStamp

	return cp, err
}

// 提交付款码支付：client.Micropay()

// 查询订单：client.QueryOrder()

// 关闭订单：client.CloseOrder()

// 撤销订单：client.Reverse()

// 申请退款：client.Refund()
func (c *WeChatPayClient) Refund(orderId string, moneyFee int64, tradeType string) bool {
	body := make(gopay.BodyMap)
	body.Set("out_trade_no", orderId)
	body.Set("nonce_str", gopay.GetRandomString(32))
	body.Set("sign_type", gopay.SignType_MD5)
	body.Set("out_refund_no", orderId)
	body.Set("total_fee", moneyFee)
	body.Set("refund_fee", moneyFee)

	//请求申请退款（沙箱环境下，证书路径参数可传空）
	//    body：参数Body
	//    certFilePath：cert证书路径
	//    keyFilePath：Key证书路径
	//    pkcs12FilePath：p12证书路径
	wxRsp, err := c.getPayClient(tradeType).Refund(body, c.getPayParam(tradeType).WeChatCertFile, c.getPayParam(tradeType).WeChatKeyFile, c.getPayParam(tradeType).WeChatP12File)
	if err != nil {
		log.Errorf("微信退款Error:%v", err)
		return false
	}
	log.Debugf("微信退款wxRsp：%#v", *wxRsp)
	if wxRsp.ReturnCode == gopay.SUCCESS {
		log.Debugf("微信退款成功:%#v", wxRsp)
		return true
	} else {
		log.Errorf("微信退款失败:%#v", wxRsp)
	}
	return false
}

// 查询退款：client.QueryRefund()
func (c *WeChatPayClient) QueryRefund(outTradeNo string, tradeType string) bool {
	//初始化参数结构体
	body := make(gopay.BodyMap)
	body.Set("out_trade_no", outTradeNo)
	//body.Set("out_refund_no", "vk4264I1UQ3Hm3E4AKsavK8npylGSgQA092f9ckUxp8A2gXmnsLEdsupURVTcaC7")
	//body.Set("transaction_id", "97HiM5j6kGmM2fk7fYMc8MgKhPnEQ5Rk")
	//body.Set("refund_id", "97HiM5j6kGmM2fk7fYMc8MgKhPnEQ5Rk")
	body.Set("nonce_str", gopay.GetRandomString(32))
	body.Set("sign_type", gopay.SignType_MD5)

	wxRsp, err := c.getPayClient(tradeType).QueryRefund(body)
	if err != nil {
		log.Errorf("查询微信退款错误:%v", err)
		return false
	}
	if wxRsp.ReturnCode == "SUCCESS" && wxRsp.RefundStatus0 == "SUCCESS" {
		return true
	}
	log.Errorf("查询微信退款返回:%+v", *wxRsp)
	return false
}

// 下载对账单：client.DownloadBill()

// 下载资金账单：client.DownloadFundFlow()

// 拉取订单评价数据：client.BatchQueryComment()

// 企业向微信用户个人付款：client.Transfer()
func (c *WeChatPayClient) Transfer(orderId, openid, userName, desc string, moneyFee int64, tradeType string) {
	nonceStr := gopay.GetRandomString(32)
	log.Debugf("partnerTradeNo:%v", orderId)
	//初始化参数结构体
	body := make(gopay.BodyMap)
	body.Set("nonce_str", nonceStr)
	body.Set("partner_trade_no", orderId)
	body.Set("openid", openid)
	body.Set("check_name", "FORCE_CHECK") // NO_CHECK：不校验真实姓名 , FORCE_CHECK：强校验真实姓名
	body.Set("re_user_name", userName)    //收款用户真实姓名。 如果check_name设置为FORCE_CHECK，则必填用户真实姓名
	body.Set("amount", moneyFee)          //企业付款金额，单位为分
	body.Set("desc", desc)                //企业付款备注，必填。注意：备注中的敏感词会被转成字符*
	body.Set("spbill_create_ip", "127.0.0.1")

	//请求申请退款（沙箱环境下，证书路径参数可传空）
	//    body：参数Body
	//    certFilePath：cert证书路径
	//    keyFilePath：Key证书路径
	//    pkcs12FilePath：p12证书路径
	wxRsp, err := c.getPayClient(tradeType).Transfer(body, c.getPayParam(tradeType).WeChatCertFile, c.getPayParam(tradeType).WeChatKeyFile, c.getPayParam(tradeType).WeChatP12File)
	if err != nil {
		log.Errorf("微信付款Error:", err)
		return
	}
	log.Debugf("wxRsp：%#v", *wxRsp)
	if wxRsp.ReturnCode == gopay.SUCCESS {
		log.Debugf("微信转账成功:%#v", wxRsp)
	} else {
		log.Error("微信转账失败:%#v", wxRsp)
	}

}

//验证微信回调
func (c *WeChatPayClient) VerifyWeChatSign(notifyReq *gopay.WeChatNotifyRequest, tradeType string) (ok bool, err error) {
	//验签操作
	return gopay.VerifyWeChatSign(c.getPayParam(tradeType).WeChatKey, gopay.SignType_MD5, notifyReq)
}
