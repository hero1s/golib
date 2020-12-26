package i18n

import (
	"errors"
	"fmt"
)

var Country = map[string]string{
	"ksa":     "ksa",
	"egypt":   "egypt",
	"qatar":   "qatar",
	"oman":    "oman",
	"kuwait":  "kuwait",
	"uae":     "uae",
	"bahrain": "bahrain",
}

var LangOrder = []string{"en", "cn"}

var Lang = map[string]int64{
	"en": 0,
	"EN": 0,
	"En": 0,
	"cn": 1,
	"CN": 1,
	"Cn": 1,
}

func GetLangOrder(lang string) int64 {
	switch lang {
	case "EN", "en", "En":
		return 0
	case "CN", "cn", "Cn":
		return 1
	default:
		// default is english
		return 0
	}
}

func WrapDatabaseError(err error) error {
	if err != nil {
		return fmt.Errorf("%v#%w", err, DatabaseError)
	}
	return err
}

/*
	客户端错误http status 400
		参数错误: 1000~1999
	服务端错误http status 500
		异常: 5000~5999
	正常逻辑错误及正确http status 200
		正常: 2000~4999
*/
const (
	ClientErrorBegin = 1000
	ClientErrorEnd   = 1999

	NormalErrorBegin = 2000
	NormalErrorEnd   = 4999

	SystemErrorBegin = 5000
	SystemErrorEnd   = 5000
)

var (
	Success = errors.New("0")
	//客户端错误
	ParamsError         = errors.New("1000")
	ParamsNotFit        = errors.New("1001")
	ParamsCannotBeEmpty = errors.New("1002")
	OperateLimit        = errors.New("1003")
	ParamsSignError     = errors.New("1004")

	//逻辑错误
	RecordNotFound = errors.New("2000")
	FrozenUser     = errors.New("2001")
	FrozenDevice   = errors.New("2002")
	CancelUser     = errors.New("2003")
	TeenageMode    = errors.New("2004")
	TeenageTimeOut = errors.New("2005")

	Unauthorized       = errors.New("4000")
	LoginFail          = errors.New("4001")
	UserAlreadyExist   = errors.New("4002")
	UserNotFound       = errors.New("4003")
	TokenInvalid       = errors.New("4004")
	TokenExpired       = errors.New("4005")
	SMSCodeError       = errors.New("4006")
	PasswordLimitError = errors.New("4007")
	FreezeLogin        = errors.New("4008")

	//服务器端错误
	SystemError     = errors.New("5000")
	DatabaseError   = errors.New("5001")
	UnknownError    = errors.New("5002")
	AccessLimit     = errors.New("5003")
	PermissionError = errors.New("5004")
)

// 错误参数定义
var ErrorCode = map[error][]string{
	Success:             {"success", "成功"},
	ParamsError:         {"params error", "参数校验错误"},
	ParamsNotFit:        {"params not fit", "参数不符合要求"},
	ParamsCannotBeEmpty: {"params cannot be empty", "参数不能为空"},
	OperateLimit:        {"operate limit", "操作频繁"},
	ParamsSignError:     {"param sign error", "参数签名错误"},

	RecordNotFound: {"record not found", "记录不存在"},
	FrozenUser:     {"frozen user", "用户已封号"},
	FrozenDevice:   {"frozen device", "该设备已被禁用"},
	CancelUser:     {"cancel user", "此账号已注销"},
	TeenageMode:    {"teenage mode", "青少年模式功能限制"},
	TeenageTimeOut: {"teenage time out", "青少年体验时间结束,请关闭青少年模式"},

	Unauthorized:       {"unauthorized access", "未授权访问"},
	LoginFail:          {"invalid username or password", "账号或密码不正确"},
	UserAlreadyExist:   {"user already exist", "用户已存在"},
	UserNotFound:       {"user not found", "用户不存在"},
	TokenInvalid:       {"token invalid", "无效的token"},
	TokenExpired:       {"token expired", "token已过期"},
	SMSCodeError:       {"smscode error", "验证码错误"},
	PasswordLimitError: {"password limit error", "密码错误次数太多,请用验证码登录"},
	FreezeLogin:        {"login freeze", "登录错误次数太多,账号已被冻结"},
	SystemError:        {"system error", "系统异常"},
	DatabaseError:      {"database error", "数据库错误"},
	UnknownError:       {"unknown error", "未知错误"},
	AccessLimit:        {"more max access", "超过并发数触发访问限速"},
	PermissionError:    {"permission error", "权限访问错误"},
}

// 获取错误信息
func GetErrorMsg(code error, langIndex int64) string {
	if code == nil {
		return ""
	}
	if v, ok := ErrorCode[code]; ok {
		return v[langIndex]
	} else {
		return code.Error()
	}
}
