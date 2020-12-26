package regex

import (
	"errors"
	"fmt"
	"git.moumentei.com/plat_go/golib/helpers/status"
	"regexp"
)

// Regular expression patterns
const (

	/** 匹配日期 */
	DateRegex = `(([0-9]{3}[1-9]|[0-9]{2}[1-9][0-9]{1}|[0-9]{1}[1-9][0-9]{2}|[1-9][0-9]{3})-(((0[13578]|1[02])-(0[1-9]|[12][0-9]|3[01]))|((0[469]|11)-(0[1-9]|[12][0-9]|30))|(02-(0[1-9]|[1][0-9]|2[0-8]))))|((([0-9]{2})(0[48]|[2468][048]|[13579][26])|((0[48]|[2468][048]|[3579][26])00))-02-29)`

	/** 匹配时间 */
	TimeRegex = `(?i)\d{1,2}:\d{2} ?(?:[ap]\.?m\.?)?|\d[ap]\.?m\.?`

	/** 匹配电话号码 */
	PhoneRegex = `^((1[3,5,8][0-9])|(14[5,7])|(17[0,6,7,8])|(19[6,7,8,9]))\d{8}$`

	/** 匹配网络连接地址 */
	UrlRegex = `^((ht|f)tps?):\\/\\/[\\w\\-]+(\\.[\\w\\-]+)+([\\w\\-\\.,@?^=%&:\\/~\\+#]*[\\w\\-\\@?^=%&\\/~\\+#])?$`

	/** 匹配邮箱地址 */
	EmailRegex = `^[a-zA-Z0-9_.-]+@[a-zA-Z0-9_-]+(.[a-zA-Z0-9_-]+)+$`

	/** 匹配Ip4 */
	IPv4Regex = `(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`

	/** 匹配Ip6 */
	Ipv6Regex = `(?:(?:(?:[0-9A-Fa-f]{1,4}:){7}(?:[0-9A-Fa-f]{1,4}|:))|(?:(?:[0-9A-Fa-f]{1,4}:){6}(?::[0-9A-Fa-f]{1,4}|(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){5}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,2})|:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){4}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,3})|(?:(?::[0-9A-Fa-f]{1,4})?:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){3}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,4})|(?:(?::[0-9A-Fa-f]{1,4}){0,2}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){2}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,5})|(?:(?::[0-9A-Fa-f]{1,4}){0,3}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){1}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,6})|(?:(?::[0-9A-Fa-f]{1,4}){0,4}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?::(?:(?:(?::[0-9A-Fa-f]{1,4}){1,7})|(?:(?::[0-9A-Fa-f]{1,4}){0,5}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(?:%.+)?\s*`

	/** 匹配Ip地址 */
	IPRegex = IPv4Regex + `|` + Ipv6Regex

	/** 匹配端口 */
	PortRegex = `6[0-5]{2}[0-3][0-5]|[1-5][\d]{4}|[2-9][\d]{3}|1[1-9][\d]{2}|10[3-9][\d]|102[4-9]`

	/** 身份证号 */
	IDCardRegex = `^[1-9]\d{5}(18|19|([23]\d))\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$`

	/** 匹配Mac地址 */
	MacRegex = `(([a-fA-F0-9]{2}[:-]){5}([a-fA-F0-9]{2}))`

	/** 匹配带表带单符号的 描述信息 */
	DescRegex = `^[，|。|！|：|“|”|？|；|《|》|（|）|\u4e00-\u9fa5A-Za-z0-9_/./,]+$`

	/** 匹配中文 */
	ChineseRegex = `^[\u4e00-\u9fa5]+$`

	/** 匹配中文名 */
	ChineseNameRegex = `^[\u4e00-\u9fa5|·]+$`

	/** 匹配英文名 */
	EnglishNameRegex = `^[a-zA-Z]+[' ']?[a-zA-Z]+$`

	/** 匹配强密码    字母、数字和符号必须同时存在，符号存在开头和结尾且仅限!@#$%^*  */
	StrongPassword = `(^[A-Za-z]+[0-9]+[!@#$%^&*]+$)|(^[0-9]+[A-Za-z]+[!@#$%^&*]+$)|(^[!@#$%^&*]+[A-Za-z]+[0-9]+$)|(^[!@#$%^&*]+[0-9]+[A-Za-z]+$)|(^[A-Za-z]+[!@#$%^&*]+[0-9]+$)|(^[0-9]+[!@#$%^&*]+[A-Za-z]+$)`

	/** 匹配中度密码  字母数字必须组合，字母数字必须同时存在*/
	MediumPassword = `(^[A-Za-z]+[0-9]+$|(^[0-9]+[A-Za-z]+$)`

	/** 匹配数字字母 */
	LetNumRegex = `^[A-Za-z0-9]+$`

	/** 匹配纯数字字符串 */
	NumStrRegex = `^[0-9]+$`

	/** 匹配十六进制字符串 */
	HexStrRegex = `^[0-9A-Fa-f]+$`
)

// 匹配是否为日期格式
func IsDate(text string) (ok bool, err error) {
	if matched, _ := regexp.MatchString(DateRegex, text); matched {
		return true, nil
	}
	return false, errors.New("请输入正确的日期格式")
}

// 匹配是否是时间格式
func IsTime(text string) (ok bool, err error) {
	if matched, _ := regexp.MatchString(TimeRegex, text); matched {
		return true, nil
	}
	return false, errors.New("请输入正确的时间格式")
}

// 匹配是否是电话号码
func IsPhone(text string) (ok bool, err error) {
	if matched, _ := regexp.MatchString(PhoneRegex, text); matched {
		return true, nil
	}
	return false, errors.New("请输入正确的手机号码")
}

// 匹配网络连接地址
func IsUrl(text string) (ok bool, err error) {
	if matched, _ := regexp.MatchString(UrlRegex, text); matched {
		return true, nil
	}
	return false, errors.New("请输入合法的url地址")
}

// 匹配邮箱地址
func IsEmail(text string) (ok bool, err error) {
	if matched, _ := regexp.MatchString(EmailRegex, text); matched {
		return true, nil
	}
	return false, errors.New("请输入正确的邮箱地址")
}

// 匹配Ip
func IsIPAddress(text string) (ok bool, err error) {
	if matched, _ := regexp.MatchString(IPRegex, text); matched {
		return true, nil
	}
	return false, errors.New("请输入合法的IP地址")
}

// 匹配 mac 地址
func IsMacAddress(text string) (ok bool, err error) {
	if matched, _ := regexp.MatchString(MacRegex, text); matched {
		return true, nil
	}
	return false, errors.New("请输入合法的Mac地址")
}

// 匹配端口
func IsPort(text string) (ok bool, err error) {
	if matched, _ := regexp.MatchString(PortRegex, text); matched {
		return true, nil
	}
	return false, errors.New("请输入合法的端口号")
}

// 匹配身份证号码
func IsIDCard(text string) (ok bool, err error) {
	if matched, _ := regexp.MatchString(IDCardRegex, text); matched {
		return true, nil
	}
	return false, errors.New("请输入正确的身份证号")
}

// 匹配带最大，最小长度限制的描述信息
func DescMatchMinAndMax(text string, min int, max int) (ok bool, err error) {
	if min > max {
		return false, errors.New("备注或描述" + status.StatusText(status.MinThanMaxErr))
	}
	if len(text) >= min && len(text) <= max {
		if matched, _ := regexp.MatchString(DescRegex, text); matched {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf("备注或描述为%d-%d位中文字符字母或数字组合", min, max))
	}
	return false, errors.New(fmt.Sprintf("备注或描述为%d-%d位中文字符字母或数字组合", min, max))
}

// 匹配带最大长度的描述信息
func DescMatchMax(text string, max int) (ok bool, err error) {
	if max < 0 {
		return false, errors.New("备注或描述" + status.StatusText(status.MaxLessZeroErr))
	}
	if len(text) <= max {
		if matched, _ := regexp.MatchString(DescRegex, text); matched {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf("备注或描述为小于%d位的中文字符字母或数字组合", max))
	}
	return false, errors.New(fmt.Sprintf("备注或描述为小于%d位的中文字符字母或数字组合", max))
}

// 匹配带最大，最小长度的数字字母
func MatchLetterNumMinAndMax(text string, min int, max int, purpose string) (ok bool, err error) {
	if min > max {
		return false, errors.New(purpose + status.StatusText(status.MinThanMaxErr))
	}
	if len(text) >= min && len(text) <= max {
		if matched, _ := regexp.MatchString(LetNumRegex, text); matched {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf(purpose+"为%d-%d位符字母或数字组合", min, max))
	}
	return false, errors.New(fmt.Sprintf(purpose+"为%d-%d位符字母或数字组合", min, max))
}

// 匹配带最大长度的数字字母
func MatchLetterNumMax(text string, max int, purpose string) (ok bool, err error) {
	if max < 0 {
		return false, errors.New(purpose + status.StatusText(status.MaxLessZeroErr))
	}
	if len(text) <= max {
		if matched, _ := regexp.MatchString(LetNumRegex, text); matched {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf(purpose+"为小于%d位符字母或数字组合", max))
	}
	return false, errors.New(fmt.Sprintf(purpose+"为小于%d位符字母或数字组合", max))
}

// 匹配带最大最小长度的中文
func MatchChineseMinAndMax(text string, min int, max int, purpose string) (ok bool, err error) {
	if min > max {
		return false, errors.New(purpose + status.StatusText(status.MinThanMaxErr))
	}
	if len(text) >= min && len(text) <= max {
		if matched, _ := regexp.MatchString(ChineseRegex, text); matched {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf(purpose+"为%d-%d位中文字符", min, max))
	}
	return false, errors.New(fmt.Sprintf(purpose+"为%d-%d位中文字符", min, max))
}

// 匹配带最大长度的中文
func MatchChineseMax(text string, max int, purpose string) (ok bool, err error) {
	if max < 0 {
		return false, errors.New(purpose + status.StatusText(status.MaxLessZeroErr))
	}
	if len(text) <= max {
		if matched, _ := regexp.MatchString(ChineseRegex, text); matched {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf(purpose+"为小于%d位中文字符", max))
	}
	return false, errors.New(fmt.Sprintf(purpose+"为小于%d位中文字符", max))
}

// 匹配中文名字
func MatchChineseName(text string, min int, max int) (ok bool, err error) {
	if min > max {
		return false, errors.New("中文名" + status.StatusText(status.MinThanMaxErr))
	}
	if len(text) >= min && len(text) <= max {
		if matched, _ := regexp.MatchString(ChineseNameRegex, text); matched {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf(status.StatusText(status.ChineseNameErr), min, max))
	}
	return false, errors.New(fmt.Sprintf(status.StatusText(status.ChineseNameErr), min, max))
}

// 匹配英文名
func MatchEnglishName(text string, min int, max int) (ok bool, err error) {
	if min > max {
		return false, errors.New("英文名" + status.StatusText(status.MinThanMaxErr))
	}
	if len(text) >= min && len(text) <= max {
		if matched, _ := regexp.MatchString(EnglishNameRegex, text); matched {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf(status.StatusText(status.EnglishNameErr), min, max))
	}
	return false, errors.New(fmt.Sprintf(status.StatusText(status.EnglishNameErr), min, max))
}

// 强密码
func MatchStrongPassword(text string, min int, max int) (ok bool, err error) {
	if min > max {
		return false, errors.New("密码" + status.StatusText(status.MinThanMaxErr))
	}
	if len(text) >= min && len(text) <= max {
		if matched, _ := regexp.MatchString(StrongPassword, text); matched {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf(status.StatusText(status.StrongPasswordErr), min, max))
	}
	return false, errors.New(fmt.Sprintf(status.StatusText(status.StrongPasswordErr), min, max))
}

// 匹配中度密码
func MatchMediumPassword(text string, min int, max int) (ok bool, err error) {
	if min > max {
		return false, errors.New("密码" + status.StatusText(status.MinThanMaxErr))
	}
	if len(text) >= min && len(text) <= max {
		if matched, _ := regexp.MatchString(MediumPassword, text); matched {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf(status.StatusText(status.MediumPasswordErr), min, max))
	}
	return false, errors.New(fmt.Sprintf(status.StatusText(status.MediumPasswordErr), min, max))
}

// 匹配纯数字字符串
func MatchNumStrMinAndMax(text string, min int, max int, purpose string) (ok bool, err error) {
	if min > max {
		return false, errors.New(purpose + status.StatusText(status.MinThanMaxErr))
	}
	if len(text) >= min && len(text) <= max {
		if matched, _ := regexp.MatchString(NumStrRegex, text); matched {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf(purpose+"为%d-%d位0-9的数字", min, max))
	}
	return false, errors.New(fmt.Sprintf(purpose+"为%d-%d位0-9的数字", min, max))
}

// 匹配定长的纯数字
func MatchNumStrFix(text string, fix int, purpose string) (ok bool, err error) {
	if fix < 0 {
		return false, errors.New(purpose + status.StatusText(status.FixLessZeroErr))
	}
	if len(text) == fix {
		if matched, _ := regexp.MatchString(HexStrRegex, text); matched {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf(purpose+"为%d位数字", fix))
	}
	return false, errors.New(fmt.Sprintf(purpose+"为%d位数字", fix))
}

// 匹配带最大最小长度的十六进制字符串
func MatchHexStrMinAndMax(text string, min int, max int) bool {
	if min > max {
		return false
	}
	if len(text) >= min && len(text) <= max {
		matched, _ := regexp.MatchString(HexStrRegex, text)
		return matched
	}
	return false
}

// 匹配带最大长度的十六进制字符串
func MatchHexStrMax(text string, max int) bool {
	if max < 0 {
		return false
	}
	if len(text) <= max {
		matched, _ := regexp.MatchString(HexStrRegex, text)
		return matched
	}
	return false
}

// 匹配定长的十六进制字符串
func MatchHexStrFixed(text string, fix int) bool {
	if fix < 0 {
		return false
	}
	if len(text) == fix {
		matched, _ := regexp.MatchString(HexStrRegex, text)
		return matched
	}
	return false
}
