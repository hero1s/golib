package sdk

// OrderBaseInfo 订单基本信息
type OrderBaseInfo struct {
	OrderID   string // 商户订单号(必填, 保持唯一性,64个英文字符以内)
	RealName  string // 收款人姓名(必填)
	IDCard    string // 收款人手机号(必填)
	Pay       string // 打款金额(必填 单位:元)
	PayRemark string // 打款备注(选填, 最大20个字符,一个汉字占两个字符,不允许特殊字符)
	NotifyURL string // 回调地址(选填, 最大长度为200)
}

// BankOrderParam 银行卡订单信息
type BankOrderParam struct {
	OrderBaseInfo
	CardNo  string // 收款人银行卡号(必填)
	PhoneNo string // 收款人手机号(选填)
}

// AliOrderParam 支付宝订单信息
type AliOrderParam struct {
	OrderBaseInfo
	CardNo    string // 收款人支付宝号(必填)
	CheckName string // 校验支付宝账户姓名(可填Check、 NoCheck)
}

// WxOrderParam 微信订单信息
type WxOrderParam struct {
	OrderBaseInfo
	OpenID    string // 商户AppID下用户OpenID(必填)
	Notes     string // 描述信息(选填)
	WxAppID   string // 微信打款商户微信AppID(选填,最⼤长度为200)
	WxPayMode string // 微信打款模式(选填，二种取值，可填"", "redpacket")
}

// BaseResponse 基础响应信息
type BaseResponse struct {
	Code      string `json:"code"`       // 响应码
	Message   string `json:"message"`    // 响应信息
	RequestID string `json:"request_id"` // 请求ID
}

// BaseCheckResponse 基础校验响应信息
type BaseCheckResponse struct {
	BaseResponse
	Data struct {
		Ok bool `json:"ok"` // 是否成功
	} `json:"data"`
}

// CreateOrderResponse 创建订单响应信息
type CreateOrderResponse struct {
	BaseResponse
	Data struct {
		Pay     string `json:"pay"`      // 打款金额
		Ref     string `json:"ref"`      // 综合服务平台订单流水号
		OrderID string `json:"order_id"` // 商户订单流水号
	} `json:"data"`
}

// QueryOrderResponse 查询订单响应信息
type QueryOrderResponse struct {
	BaseResponse
	Data OrderInfo `json:"data"`
}

// OrderInfo 订单详细信息
type OrderInfo struct {
	AnchorID            string `json:"anchor_id"`
	BrokerAmount        string `json:"broker_amount"`
	BrokerBankBill      string `json:"broker_bank_bill"`
	BrokerFee           string `json:"broker_fee"`
	BrokerID            string `json:"broker_id"`
	CardNo              string `json:"card_no"`
	CreatedAt           string `json:"created_at"`
	DealerID            string `json:"dealer_id"`
	EncryData           string `json:"encry_data"`
	FeeAmount           string `json:"fee_amount"`
	FinishedTime        string `json:"finished_time"`
	IDCard              string `json:"id_card"`
	Notes               string `json:"notes"`
	OrderID             string `json:"order_id"`
	Pay                 string `json:"pay"`
	PayRemark           string `json:"pay_remark"`
	PhoneNo             string `json:"phone_no"`
	RealName            string `json:"real_name"`
	Ref                 string `json:"ref"`
	Status              string `json:"status"`
	StatusDetail        string `json:"status_detail"`
	StatusDetailMessage string `json:"status_detail_message"`
	StatusMessage       string `json:"status_message"`
	SysAmount           string `json:"sys_amount"`
}

// OrderDetailInfo 回调通知订单信息
type OrderDetailInfo struct {
	AnchorID            string `json:"anchor_id"`
	BrokerAmount        string `json:"broker_amount"`
	BrokerBankBill      string `json:"broker_bank_bill"`
	BrokerFee           string `json:"broker_fee"`
	BrokerID            string `json:"broker_id"`
	BrokerWalletRef     string `json:"broker_wallet_ref"`
	CardNo              string `json:"card_no"`
	CreatedAt           string `json:"created_at"`
	DealerID            string `json:"dealer_id"`
	FinishedTime        string `json:"finished_time"`
	IDCard              string `json:"id_card"`
	Notes               string `json:"notes"`
	OrderID             string `json:"order_id"`
	Pay                 string `json:"pay"`
	PayRemark           string `json:"pay_remark"`
	PhoneNo             string `json:"phone_no"`
	RealName            string `json:"real_name"`
	Ref                 string `json:"ref"`
	Status              string `json:"status"`
	StatusDetail        string `json:"status_detail"`
	StatusDetailMessage string `json:"status_detail_message"`
	StatusMessage       string `json:"status_message"`
	SysAmount           string `json:"sys_amount"`
	SysBankBill         string `json:"sys_bank_bill"`
	SysFee              string `json:"sys_fee"`
	SysWalletRef        string `json:"sys_wallet_ref"`
	Tax                 string `json:"tax"`
	UserFee             string `json:"user_fee"`
	WithdrawPlatform    string `json:"withdraw_platform"`
}

// OrderCallBackResponse 订单回调信息
type OrderCallBackResponse struct {
	NotifyID   string          `json:"notify_id"`
	NotifyTime string          `json:"notify_time"`
	Data       OrderDetailInfo `json:"data"`
}

// QueryAccountBalanceResponse 查询商户余额响应信息
type QueryAccountBalanceResponse struct {
	BaseResponse
	Data struct {
		DealerInfos []AccountBalance `json:"dealer_infos"`
	} `json:"data"`
}

// AccountBalance 账户余额信息
type AccountBalance struct {
	BrokerID         string `json:"broker_id"`          // 代征主体ID
	BankCardBalance  string `json:"bank_card_balance"`  // 银行卡余额
	AlipayBalance    string `json:"alipay_balance"`     // ⽀付宝余额
	WxpayBalance     string `json:"wxpay_balance"`      // 微信余额
	IsBankCard       bool   `json:"is_bank_card"`       // 是否开通银行卡通道
	IsAlipay         bool   `json:"is_alipay"`          // 是否开通付宝通道
	IsWxpay          bool   `json:"is_wxpay"`           // 是否开通微信通道
	RebateFeeBalance string `json:"rebate_fee_balance"` // 服务费返点余额
}

// QueryReceiptFileResponse 查询电子回单响应信息
type QueryReceiptFileResponse struct {
	BaseResponse
	Data OrderReceiptFile `json:"data"`
}

// OrderReceiptFile 电子回单信息
type OrderReceiptFile struct {
	ExpireTime string `json:"expire_time"` // 过期时间
	FileName   string `json:"file_name"`   // 文件名称
	URL        string `json:"url"`         // 下载地址
}

// DownloadOrderFileResponse 下载日订单响应信息
type DownloadOrderFileResponse struct {
	BaseResponse
	Data struct {
		OrderDownloadURL string `json:"order_download_url"` // url地址
	} `json:"data"`
}

// DownloadBillFileResponse 下载日流水响应信息
type DownloadBillFileResponse struct {
	BaseResponse
	Data struct {
		BillDownloadURL string `json:"bill_download_url"` // url地址
	} `json:"data"`
}

// QueryRechargeRecordResponse 充值记录响应信息
type QueryRechargeRecordResponse struct {
	BaseResponse
	Data []RechargeRecord `json:"data"`
}

// RechargeRecord 充值记录信息
type RechargeRecord struct {
	BrokerID        string `json:"broker_id"`        // 代征主体ID
	DealerID        string `json:"dealer_id"`        // 商户ID
	ActualAmount    int    `json:"actual_amount"`    // 实际到账金额
	Amount          int    `json:"amount"`           // 充值金额
	CreatedAt       string `json:"created_at"`       // 创建时间
	RechargeChannel string `json:"recharge_channel"` // 充值渠道
	RechargeID      string `json:"recharge_id"`      // 充值记录ID
}

// UserInfoParam 免验证用户名单信息
type UserInfoParam struct {
	RealName     string   `json:"real_name"`     // 姓名
	IDCard       string   `json:"id_card"`       // 证件号
	Birthday     string   `json:"birthday"`      // 出生日期
	CardType     string   `json:"card_type"`     // 证件类型
	Country      string   `json:"country"`       // 国家代码
	Gender       string   `json:"gender"`        // 性别
	NotifyURL    string   `json:"notify_url"`    // 回调地址
	Ref          string   `json:"ref"`           // 流水号(回调时附带)
	UserImages   []string `json:"user_images"`   // 证件照片
	CommentApply string   `json:"comment_apply"` // 申请备注
}

// UserCallBackInfo 通知用户上传信息
type UserCallBackInfo struct {
	BrokerID string `json:"broker_id"` // 代征主体ID
	DealerID string `json:"dealer_id"` // 商户ID
	Comment  string `json:"comment"`   // 备注
	RealName string `json:"real_name"` // 姓名
	IDCard   string `json:"id_card"`   // 证件号
	Ref      string `json:"ref"`       // 凭证(上传信息中)
	Status   string `json:"status"`    // 状态(pass: 通过 reject: 拒绝)
}

// QueryInvoiceResponse 查询发票响应信息
type QueryInvoiceResponse struct {
	BaseResponse
	Data InvoiceInfo `json:"data"`
}

// InvoiceInfo 发票信息
type InvoiceInfo struct {
	BrokerID    string `json:"broker_id"`    // 代征主体ID
	DealerID    string `json:"dealer_id"`    // 商户ID
	Invoiced    string `json:"invoiced"`     // 已开发票金额
	NotInvoiced string `json:"not_invoiced"` // 待开发票⾦额
}

// ElementVerifyResponse 银行卡四要素发送短信请求信息
type ElementVerifyResponse struct {
	BaseResponse
	Data struct {
		Ref string `json:"ref"` // 交易凭证
	} `json:"data"`
}

// QueryBankCardResponse 查询银行卡信息响应信息
type QueryBankCardResponse struct {
	BaseResponse
	Data BankCardInfo `json:"data"` // 银行卡信息
}

// BankCardInfo 银行卡信息
type BankCardInfo struct {
	BankCode  string `json:"bank_code"`  // 银行代码
	BankName  string `json:"bank_name"`  // 银行名称
	CardType  string `json:"card_type"`  // 银行卡类型
	IsSupport bool   `json:"is_support"` // 云账户综合服务平台是否支持该银行打款
}
