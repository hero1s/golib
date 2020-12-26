package gopay

import (
	"git.moumentei.com/plat_go/golib/log"
	"github.com/iGoogle-ink/gopay"
)

type AliPayClient struct {
	PayClient *gopay.AliPayClient
	PayCfg    AliPayParam `json:"ali_pay" desc:"支付包支付参数"`
}

func InitAliPay(isProd bool, payCfg AliPayParam) *AliPayClient {
	client := new(AliPayClient)
	client.PayCfg = payCfg
	client.PayClient = gopay.NewAliPayClient(payCfg.AliAppId, string(keyFromFile(payCfg.AliAppPrivateKeyFile)), isProd)

	//设置支付宝请求 公共参数
	client.PayClient.SetCharset("utf-8").
		SetSignType("RSA2"). //设置签名类型，不设置默认 RSA2
		SetReturnUrl(payCfg.AliReturnUrl). //设置返回URL
		SetNotifyUrl(payCfg.AliAppCallbackUrl) //设置异步通知URL
	//.SetAppAuthToken().SetAuthToken()
	return client
}

//* 手机网站支付接口2.0（手机网站支付）：client.AliPayTradeWapPay()
func (c *AliPayClient) AliPayTradeWapPay(moneyFee int64, describe, orderId, quitUrl, returnUrl string,timeout string) (string, error) {
	if len(returnUrl) > 1 {
		c.PayClient.SetReturnUrl(returnUrl)
	}
	//请求参数
	body := make(gopay.BodyMap)
	body.Set("subject", describe)
	body.Set("out_trade_no", orderId)
	body.Set("quit_url", quitUrl)
	body.Set("total_amount", AliMoneyFeeToString(float64(moneyFee)))
	body.Set("product_code", "QUICK_WAP_WAY")
	body.Set("timeout_express",timeout)
	//手机网站支付请求
	payUrl, err := c.PayClient.AliPayTradeWapPay(body)
	if err != nil {
		log.Errorf("page pay err:", err)
		return payUrl, err
	}
	return payUrl, err
}

//* 统一收单下单并支付页面接口（电脑网站支付）：client.AliPayTradePagePay()
func (c *AliPayClient) AliPayTradePagePay(moneyFee int64, describe, orderId string, returnUrl string,timeout string) (string, error) {
	if len(returnUrl) > 1 {
		c.PayClient.SetReturnUrl(returnUrl)
	}
	//请求参数
	body := make(gopay.BodyMap)
	body.Set("subject", describe)
	body.Set("out_trade_no", orderId)
	body.Set("total_amount", AliMoneyFeeToString(float64(moneyFee)))
	body.Set("product_code", "FAST_INSTANT_TRADE_PAY")
	body.Set("timeout_express",timeout)
	//电脑网站支付请求
	payUrl, err := c.PayClient.AliPayTradePagePay(body)
	if err != nil {
		log.Errorf("page pay err:", err)
		return payUrl, err
	}
	return payUrl, err
}

//* APP支付接口2.0（APP支付）：client.AliPayTradeAppPay()
func (c *AliPayClient) AliPayTradeAppPay(moneyFee int64, describe, orderId string,timeout string) (string, error) {
	//请求参数
	body := make(gopay.BodyMap)
	body.Set("subject", describe)
	body.Set("out_trade_no", orderId)
	body.Set("total_amount", AliMoneyFeeToString(float64(moneyFee)))
	body.Set("product_code", c.PayCfg.AliProductCode)
	body.Set("timeout_express",timeout)
	//手机APP支付参数请求
	payParam, err := c.PayClient.AliPayTradeAppPay(body)
	if err != nil {
		log.Errorf("阿里支付预下单失败err:", err)
	}
	//log.Infof("阿里支付预下单返回payParam:", payParam)
	return payParam, err
}

//* 统一收单交易支付接口（商家扫用户付款码）：client.AliPayTradePay()

//* 统一收单交易创建接口（小程序支付）：client.AliPayTradeCreate()

//* 统一收单线下交易查询：client.AliPayTradeQuery()

//* 统一收单交易关闭接口：client.AliPayTradeClose()

//* 统一收单交易撤销接口：client.AliPayTradeCancel()

//* 统一收单交易退款接口：client.AliPayTradeRefund()
func (c *AliPayClient) AliPayTradeRefund(orderId string, moneyFee int64) bool {
	//请求参数
	body := make(gopay.BodyMap)
	body.Set("out_trade_no", orderId)
	body.Set("refund_amount", AliMoneyFeeToString(float64(moneyFee)))
	body.Set("refund_reason", "测试退款")
	//发起退款请求
	aliRsp, err := c.PayClient.AliPayTradeRefund(body)
	if err != nil {
		log.Errorf("阿里退款失败err:%v", err)
		return false
	}
	log.Infof("阿里退款返回aliRsp:%#v", *aliRsp)
	if aliRsp.AlipayTradeRefundResponse.Code == "10000" {
		log.Debugf("阿里退款成功:%#v", aliRsp)
		return true
	} else {
		log.Debugf("阿里退款失败:%#v", aliRsp)
	}
	return false
}

//* 统一收单退款页面接口：client.AliPayTradePageRefund()

//* 统一收单交易退款查询：client.AliPayTradeFastPayRefundQuery()
func (c *AliPayClient) AliPayTradeFastPayRefundQuery(outTradeNo string) string {
	//请求参数
	body := make(gopay.BodyMap)
	body.Set("out_trade_no", outTradeNo)
	body.Set("out_request_no", outTradeNo)
	//发起退款查询请求
	aliRsp, err := c.PayClient.AliPayTradeFastPayRefundQuery(body)
	if err != nil {
		log.Errorf("支付宝退款查询错误:%v", err)
		return "0"
	}
	if aliRsp.AliPayTradeFastpayRefundQueryResponse.Code != "10000" {
		log.Error("阿里退款查询返回:%+v", *aliRsp)
		return "0"
	}
	return aliRsp.AliPayTradeFastpayRefundQueryResponse.RefundAmount
}

//* 统一收单交易结算接口：client.AliPayTradeOrderSettle()

//* 统一收单线下交易预创建（用户扫商品收款码）：client.AliPayTradePrecreate()

//* 单笔转账到支付宝账户接口（商户给支付宝用户转账）：client.AlipayFundTransToaccountTransfer()
func (c *AliPayClient) AlipayFundTransToaccountTransfer(account string, moneyFee int64, desc string) {
	body := make(gopay.BodyMap)
	out_biz_no := gopay.GetRandomString(32)
	body.Set("out_biz_no", out_biz_no)
	body.Set("payee_type", "ALIPAY_LOGONID")
	body.Set("payee_account", account)
	body.Set("amount", AliMoneyFeeToString(float64(moneyFee)))
	body.Set("payer_show_name", "发钱人名字")
	body.Set("payee_real_name", "收钱人名字")
	body.Set("remark", desc)
	//创建订单
	aliRsp, err := c.PayClient.AlipayFundTransToaccountTransfer(body)
	if err != nil {
		log.Error("阿里转账支付错误err:", err)
		return
	}
	log.Debugf("阿里转账支付返回aliRsp:%#v", *aliRsp)
	if aliRsp.AlipayFundTransToaccountTransferResponse.Code == "10000" {
		log.Debugf("支付宝转账成功:%#v", aliRsp)
	} else {
		log.Error("支付宝转账失败:%#v", aliRsp)
	}
}

//* 换取授权访问令牌（获取access_token，user_id等信息）：client.AliPaySystemOauthToken()

//* 支付宝会员授权信息查询接口（App支付宝登录）：client.AlipayUserInfoShare()

//* 换取应用授权令牌（获取app_auth_token，auth_app_id，user_id等信息）：client.AlipayOpenAuthTokenApp()

//* 获取芝麻信用分：client.ZhimaCreditScoreGet()

//验证支付宝回调
func (c *AliPayClient) VerifyAliPaySign(notifyReq *gopay.AliPayNotifyRequest) (ok bool, err error) {
	return gopay.VerifyAliPaySign(string(keyFromFile(c.PayCfg.AliPayPublicKeyFile)), notifyReq)
}
