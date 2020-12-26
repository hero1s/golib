package sdk

// 基础配置信息
const (
	BrokerID = "yiyun73"  // 代征主体ID
	DealerID = "05476996" // 商户ID
)

// 基础订单状态
const (
	OrderDelete    = -1 // 订单删除
	OrderAccept    = 0  // 订单已受理
	OrderSuccess   = 1  // 订单已打款
	OrderFailed    = 2  // 订单已失败
	OrderPending   = 4  // 订单待打款(暂停处理)
	OrderSending   = 5  // 订单打款中
	OrderReadySend = 8  // 订单待打款
	OrderReturned  = 9  // 订单已退汇
	OrderCancel    = 15 // 订单取消
)

// 路由信息
const (
	BaseURL = "https://api-jiesuan.yunzhanghu.com" // 基础url

	BankOrderURL        = "/api/payment/v1/order-realtime"               // 银行卡下单接口url
	AliOrderURL         = "/api/payment/v1/order-alipay"                 // 支付宝下单接口url
	WxOrderURL          = "/api/payment/v1/order-wxpay"                  // 微信下单接口url
	QueryOrderURL       = "/api/payment/v1/query-realtime-order"         // 查单接口url
	CancelOrderURL      = "/api/payment/v1/order/fail"                   // 取消订单url
	QueryAccountURL     = "/api/payment/v1/query-accounts"               // 查询账户信息url
	QueryReceiptFileURL = "/api/payment/v1/receipt/file"                 // 查询电子回单URL
	QueryRechargeURL    = "/api/dataservice/v2/recharge-record"          // 查询充值记录url
	DownloadOrderURL    = "/api/dataservice/v1/order/downloadurl"        // 下载日订单url
	DownloadBillURL     = "/api/dataservice/v2/bill/downloadurl"         // 下载日流水url
	UploadUserURL       = "/api/payment/v1/user/exempted/info"           // 上传用户免验证名单url
	CheckExistUserURL   = "/api/payment/v1/user/white/check"             // 校验免验证用户名单是否存在url
	QueryInvoiceURL     = "/api/payment/v1/invoice-stat"                 // 查询发票接口
	Element4RequestURL  = "/authentication/verify-request"               // 银行卡四要素鉴权发送短信url
	Element4ConfirmURL  = "/authentication/verify-confirm"               // 银行卡四要素鉴权提交验证码url
	Element4URL         = "/authentication/verify-bankcard-four-factor"  // 银行卡四要素鉴权url
	Element3URL         = "/authentication/verify-bankcard-three-factor" // 银行卡三要素鉴权url
	IDCheckURL          = "/authentication/verify-id"                    // 实名制二要素鉴权url
	BankCardInfoURL     = "/api/payment/v1/card"                         // 银行卡信息查询url
)

// 状态标识码
const (
	SuccessCode = "0000" // 成功
)
