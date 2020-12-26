package gopay

import (
	"github.com/shopspring/decimal"
	"io/ioutil"
	"strconv"
)

type WeChatPayParam struct {
	WeChatAppId       string `json:"we_chat_app_id" desc:"微信appID"`
	WeChatMchId       string `json:"we_chat_mch_id" desc:"微信支付商户号"`
	WeChatKey         string `json:"we_chat_key" desc:"微信API密钥"`
	WeChatQueryUrl    string `json:"we_chat_query_url" desc:"微信订单查询url"`
	WeChatCallbackUrl string `json:"we_chat_callback_url" desc:"支付回调接口"`
	WeChatCertFile    string `json:"we_chat_cert_file" desc:"cert证书"`
	WeChatKeyFile     string `json:"we_chat_key_file" desc:"key证书"`
	WeChatP12File     string `json:"we_chat_p12_file" desc:"p12证书路径"`
}
type AliPayParam struct {
	AliProductCode       string `json:"ali_product_code" desc:"产品码"`
	AliAppId             string `json:"ali_app_id" desc:"应用ID"`
	AliPartnerId         string `json:"ali_partner_id" desc:"支付宝合作身份ID"`
	AliSellerId          string `json:"ali_seller_id" desc:"卖家支付宝用户号"`
	AliAppCallbackUrl    string `json:"ali_app_callback_url" desc:"阿里支付回调"`
	AliReturnUrl         string `json:"ali_return_url" desc:""`
	AliPayPublicKeyFile  string `json:"ali_pay_public_key_file" desc:"支付宝公钥"`
	AliAppPrivateKeyFile string `json:"ali_app_private_key_file" desc:"应用私钥"`
}

func keyFromFile(fileName string) []byte {
	data, _ := ioutil.ReadFile(fileName)
	return data
}

// 支付宝金额转字符串(单位分-->浮点元)
func AliMoneyFeeToString(moneyFee float64) string {
	moneyFee = moneyFee / 100
	return decimal.NewFromFloat(moneyFee).Truncate(2).String()
}
func AliStringToMoneyFee(fee string) int64 {
	totalFee, _ := strconv.ParseFloat(fee, 64)
	totalFee = totalFee * 100
	return int64(totalFee)
}
