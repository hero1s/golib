package status

// 定义一些系统常用的 错误码

const (
	BindModelErr = 20200
	NoneParamErr = 20201

	LoginStatusSQLErr = 20319
	LoginStatusErr    = 20300
	LoginStatusOK     = 20301

	SaveStatusOK  = 20400
	SaveStatusErr = 20401
	SaveObjIsNil  = 20402

	DeleteStatusOK  = 20403
	DeleteStatusErr = 20404
	DeleteObjIsNil  = 20405

	UpdateObjIsNil = 20406

	ExistSameNameErr  = 20501
	ExistSamePhoneErr = 20502

	FixLessZeroErr = 20797
	MaxLessZeroErr = 20798
	MinThanMaxErr  = 20799

	MediumPasswordErr = 20801
	StrongPasswordErr = 20802
	ChineseNameErr    = 20803
	EnglishNameErr    = 20804
)

var statusText = map[int32]string{

	BindModelErr: "模型封装异常！",
	NoneParamErr: "无有效参数",

	LoginStatusSQLErr: "用户登陆时更新登陆数据异常！",
	LoginStatusErr:    "用户名或密码错误!",
	LoginStatusOK:     "登陆成功！",

	SaveStatusOK:  "保存成功！",
	SaveStatusErr: "保存失败！",
	SaveObjIsNil:  "保存的对象为空！",

	DeleteStatusOK:  "删除成功！",
	DeleteStatusErr: "删除失败！",
	DeleteObjIsNil:  "删除的记录不存在！",

	UpdateObjIsNil: "修改的记录不存在！",

	ExistSameNameErr:  "已存在同名记录！",
	ExistSamePhoneErr: "已存在相同手机号！",

	/** 与正则相关的一些错误信息 */
	FixLessZeroErr: "验证规则错误，定长小于0",
	MaxLessZeroErr: "验证规则错误，最大值小于0",
	MinThanMaxErr:  "验证规则错误，最大值小于最小值",

	MediumPasswordErr: "密码为%d-%d位字母、数字，字母数字必须同时存在",
	StrongPasswordErr: "密码为%d-%d位字母、数字和符号必须同时存在，符号存在开头和结尾且仅限!@#$%^*",
	ChineseNameErr:    "中文名为%d-%d位中文字符可包含'·'",
	EnglishNameErr:    "英文名为%d-%d英文字符可包含空格",
}

func StatusText(code int32) string {
	return statusText[code]
}
