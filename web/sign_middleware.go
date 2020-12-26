package web

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/hero1s/golib/helpers/crypto"
	"github.com/hero1s/golib/i18n"
	"github.com/hero1s/golib/log"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

//数组类型为""
func convertVtoStr(value interface{}) (string, error) {
	switch value.(type) {
	case string:
		result, _ := value.(string)
		return result, nil
	case bool:
		temp := reflect.ValueOf(value).Bool()
		return strconv.FormatBool(temp), nil
	case int, int8, int16, int32, int64:
		temp := reflect.ValueOf(value).Int()
		return strconv.FormatInt(temp, 10), nil
	case uint, uint8, uint16, uint32, uint64:
		temp := reflect.ValueOf(value).Uint()
		return strconv.FormatUint(temp, 10), nil
	case float32:
		temp := reflect.ValueOf(value).Float()
		return strconv.FormatFloat(temp, 'f', -1, 32), nil
	case float64:
		temp := reflect.ValueOf(value).Float()
		return strconv.FormatFloat(temp, 'f', -1, 64), nil
	default:
		return "", errors.New("invalid value type")
	}
}

func changeMd5Str(p map[string]interface{}) (string, error) {
	keys := make([]string, 0)
	for k := range p {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var calStr string
	for _, k := range keys {
		calStr += k
		v, err := convertVtoStr(p[k])
		if err != nil {
			log.Errorf("convert value error:%v,%v", p[k], err)
			return "", errors.New(fmt.Sprintf("参数 %v 类型错误", k))
		}
		calStr += v
	}
	return calStr, nil
}

/*
	securet:md5签名秘钥
	headKeys:head里面需要校验的字段
	1:sign 字段为签名字段,ts 为时间戳,放head
	2:取出head 指定字段headKeys的参数以及get参数,存入map[string]string
	3:对map的key排序后拼接字符串str += key + values
	4:sign = md5(str + bodystr + securet + ts)
*/
func CheckParam(securet string, headKeys []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cliSign := c.Request.Header.Get("Sign")
		if cliSign == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  -1,
				"message": i18n.ParamsSignError.Error(),
			})
			c.Abort()
			return
		}
		ts := c.Request.Header.Get("Ts")
		params := make(map[string]interface{})
		//body 参数转json
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}
		// 读取后写回
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		//json.Unmarshal(bodyBytes,&params)
		//添加head参数
		for _, k := range headKeys {
			params[k] = c.GetHeader(k)
		}
		//添加GET参数
		querys := c.Request.URL.Query()
		for k := range querys {
			params[k] = strings.Join(querys[k], "")
		}
		pstr, err := changeMd5Str(params)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  -1,
				"message": i18n.ParamsSignError.Error(),
			})
			c.Abort()
			return
		}
		pstr += string(bodyBytes) + securet + ts
		//MD5
		hashStr := crypto.Md5(pstr)
		if hashStr != cliSign {
			log.Errorf("sign:%v error:cli(%v)-->(%v)", pstr, cliSign, hashStr)
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  -1,
				"message": i18n.ParamsSignError.Error(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

/*HTTP头

Sv             :<int> 标识签名的白名单Header字段列表
Sign           :<string> 接口参数签名, 见 `口签名规则`
Ts             :<int> 时间戳
Token          :<string> 用户token（注：未登录用户的请求，Token值传空字符串""）
Uid            :<string> 用户id（注：未登录用户的请求，Uid值传空字符串""）
Uuid           :<string> 设备号
Version        :<string> APP版本号：`android_0.10.0`, `ios_0.10.0`
Device         :<string> iPhone10,6（手机型号）+ "|" + ios_12.2.2（系统版本）, 例：androd  SAMSUNG/SM-A300FU|6.0.1 ; ios iPhone7,2|ios_12.4.0
Net            :<string> 网络环境: 0（移动网络），1（wifi网络）
Ch             :<int> 渠道号
Lang           :<string> 客户端使用语言：cn(中文), en（英语), ar（阿拉伯语）, fr（法语）, tr（土耳其语）
Idfa           :<string> idfa，iOS特有
Idfv           :<string> idfa，iOS特有
Locale         :<string> 地区语言: zh-cn (转小写)
Tz             :<string> 时区，有效值：-12~12。+8代表东8区，-8代表西八区，0代表标准时区
Gad            :<string> gps_adid的base64编码, android特有
Country        :<string> xx-yy xx: sms卡归属地 android独有 iso3166-1 alpha2，yy: 手机归属地 使用协议iso3166-1 ios: alpha2 android: alpha3  相关说明：https://baike.baidu.com/item/ISO%203166-1/5269555?fr=aladdin
*/
