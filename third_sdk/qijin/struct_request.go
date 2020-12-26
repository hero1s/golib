package qijin

//打款状态
const (
	OrderStatus_UNPAID  = "UNPAID"  //待支付
	OrderStatus_WAIT    = "WAIT"    //打款中
	OrderStatus_SUCCESS = "SUCCESS" //已打款
	OrderStatus_FAIL    = "FAIL"    //已无效
	OrderStatus_FAILED  = "FAILED"  //打款失败
)

//创建用户
type UserCreation struct {
	Name   string `json:"name" desc:"实名信息"`
	IdCard string `json:"idCard" desc:"身份证号"`
	Mobile string `json:"mobile" desc:"绑卡手机号"`
	/*	IdCardFrontFileUrl string `json:"idCardFrontFileUrl,omitempty" desc:"身份证正面图片(可选)"`
		IdCardBackFileUrl  string `json:"idCardBackFileUrl,omitempty" desc:"身份证反面图片(可选)"`
		ContractFileUrl    string `json:"contractFileUrl,omitempty" desc:"签约合同文件链接(可选)"`
		LivingFileUrl      string `json:"livingFileUrl,omitempty" desc:"活体认证图片(可选)"`*/
}

type UserCreationContent struct {
	IdCard         string `json:"idCard" desc:"身份证号"`
	UserAuthStatus int64  `json:"userAuthStatus" desc:"用户认证状态(0失败1成功)"`
	Message        string `json:"message" desc:"提示信息"`
}

//创建用户返回
type UserCreationResponse struct {
	RetCode    int64               `json:"retCode" desc:"返回代码"`
	RetMessage string              `json:"retMessage" desc:"返回信息"`
	Content    UserCreationContent `json:"content,omitempty" desc:"内容"`
}

//打款订单
type OrderCreation struct {
	RequesterOrderNumber string `json:"requesterOrderNumber" desc:"请求方唯一订单号"`
	Name                 string `json:"name" desc:"实名信息"`
	IdCard               string `json:"idCard" desc:"身份证号"`
	Amount               string `json:"amount" desc:"金额"`
	BankCard             string `json:"bankCard" desc:"打款银行卡号"`
	PaymentNote          string `json:"paymentNote" desc:"打款备注"`
}

type OrderCreationContent struct {
	RequesterOrderNumber string  `json:"requesterOrderNumber" desc:"请求方唯一订单号"`
	SystemOrderNumber    string  `json:"systemOrderNumber" desc:"企进订单号"`
	Amount               float64 `json:"amount" desc:"金额"`
	ServiceCharge        float64 `json:"serviceCharge" desc:"服务费"`
	OrderStatus          string  `json:"orderStatus" desc:"订单状态"`
	VerifyCode           int64   `json:"verifyCode" desc:"校验码"`
	Message              string  `json:"message" desc:"提示信息"`
}

//打款订单返回
type OrderCreationResponse struct {
	RetCode    int64                `json:"retCode" desc:"返回代码"`
	RetMessage string               `json:"retMessage" desc:"返回信息"`
	Content    OrderCreationContent `json:"content,omitempty" desc:"内容"`
}

//打款订单查询
type OrderDetail struct {
	//SystemOrderNumber    string `json:"systemOrderNumber" desc:"打款订单号"`
	RequesterOrderNumber string `json:"requesterOrderNumber" desc:"请求方订单号"`
}
type OrderDetailContent struct {
	RequesterOrderNumber string  `json:"requesterOrderNumber" desc:"请求方唯一订单号"`
	SystemOrderNumber    string  `json:"systemOrderNumber" desc:"企进订单号"`
	Name                 string  `json:"name" desc:"实名信息"`
	IdCard               string  `json:"idCard" desc:"身份证号"`
	Amount               float64 `json:"amount" desc:"金额"`
	BankCard             string  `json:"bankCard" desc:"打款银行卡号"`
	OrderStatus          string  `json:"orderStatus" desc:"订单状态"`
	BatchNumber          string  `json:"batchNumber" desc:"批次号"`
	ServiceCharge        float64 `json:"serviceCharge" desc:"服务费"`
	CapitalFlowNumber    string  `json:"capitalFlowNumber" desc:"银行订单流水号"`
	Remark               string  `json:"remark" desc:"打款备注"`
	Reason               string  `json:"reason" desc:"打款失败原因"`
	PaymentTime          int64   `json:"paymentTime" desc:"支付时间"`
	GmtCreate            int64   `json:"gmtCreate" desc:"创建时间"`
	Message              string  `json:"message" desc:"提示信息"`
}

type OrderDetailResponse struct {
	RetCode    int64              `json:"retCode" desc:"返回代码"`
	RetMessage string             `json:"retMessage" desc:"返回信息"`
	Content    OrderDetailContent `json:"content,omitempty" desc:"内容"`
}
