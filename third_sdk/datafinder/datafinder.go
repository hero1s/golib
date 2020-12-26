package datafinder

import (
	"encoding/json"
	"errors"
	"git.moumentei.com/plat_go/golib/log"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	DataRangersSinglePostUrl = "https://mcs.snssdk.com/v2/event/json"

	ResponseCodeOk        = 200 // 成功
	ResponseCodeMalformed = 400 // 请求格式错误, 查看X-MCS-AppKey与header，user的定义
	ResponseCodeTooMany   = 413 // 单个请求事件数过多,或请求json数组元素过多（只针对list接口）
	ResponseCodeFrequency = 429 // 设备请求过于频繁
	ResponseCodeServerErr = 503 // 服务器内部错误
)

var (
	RequestParamsValidErr = errors.New("请求参数错误")
	RequestMalformedErr   = errors.New("请求格式错误")
	RequestTooManyErr     = errors.New("请求事件数过多")
	RequestFrequencyErr   = errors.New("设备请求过于频繁")
	RequestServerErr      = errors.New("请求服务错误")
)

var netTransport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: 5 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 5 * time.Second,
}
var DefaultNetClient = &http.Client{
	Timeout:   time.Second * 10,
	Transport: netTransport,
}

type DataFinder struct {
	AppKey string //应用的APP Key
	Body   *DataFinderBody
}

type DataFinderBody struct {
	User   *UserStruct    `json:"user"`
	Header *HeaderStruct  `json:"header"`
	Events []*EventStruct `json:"events"`
}

type ResponseBody struct {
	Message string `json:"message"`
	Sc      int    `json:"sc"`
}

type UserStruct struct {
	UserUniqueId string `json:"user_unique_id"` //用户的唯一身份标识， 必选
}

type HeaderStruct struct {
	AppName      string            `json:"app_name"`                 // 应用的英文名称, 必选
	AppPackage   string            `json:"app_package,omitempty"`    // 包名
	AppChannel   string            `json:"app_channel,omitempty"`    // app分发渠道
	AppVersion   string            `json:"app_version,omitempty"`    // app版本，三段分隔，如1.0.1
	OsName       string            `json:"os_name,omitempty"`        // 客户端系统
	OsVersion    string            `json:"os_version,omitempty"`     // 客户端系统版本号
	DeviceModel  string            `json:"device_model,omitempty"`   // 设备型号
	AbSdkVersion string            `json:"ab_sdk_version,omitempty"` // ab实验分组信息
	TrafficType  string            `json:"traffic_type,omitempty"`   // 流量类型
	ClientIp     string            `json:"client_ip,omitempty"`      // 客户端ip
	Custom       map[string]string `json:"custom,omitempty"`         // 自定义字段,单层json map。上述字段都是保留字段不能使用
	/* Custom自定义字段
	region			string		否		所在区域国家(系统设置)，us等，(放在custom中)
	language		string		否		语言(系统设置)，en等，(放在custom中)
	app_region		string		否		国家(app设置)，us等，(放在custom中)
	app_language 	string		否		语言(app设置)，en等，(放在custom中)
	timezone		string		否		时区，-12~12，(放在custom中)
	*/
	UtmSource   string `json:"utm_source,omitempty"`   // 推广来源
	UtmCampaign string `json:"utm_campaign,omitempty"` // 推广活动
	UtmMedium   string `json:"utm_medium,omitempty"`   // 推广媒介
	UtmContent  string `json:"utm_content,omitempty"`  // 推广内容
	UtmTerm     string `json:"utm_term,omitempty"`     // 推广关键词
}

type EventStruct struct {
	EventName   string `json:"event"`         // 事件名, 必选
	Params      string `json:"params"`        // 事件参数,单层json map, 必选
	LocalTimeMs int64  `json:"local_time_ms"` // unix_timestamp( 毫秒), 必选
}

func NewDataFinder(appKey string) *DataFinder {
	return &DataFinder{
		AppKey: appKey,
		Body:   &DataFinderBody{},
	}
}

func (df *DataFinder) isInit() bool {
	if df.AppKey == "" || df.Body == nil {
		return false
	}
	return true
}

func (df *DataFinder) valid() bool {
	if !df.isInit() {
		return false
	}

	if df.Body.User == nil || df.Body.User.UserUniqueId == "" {
		return false
	}
	if df.Body.Header == nil || df.Body.Header.AppName == "" {
		return false
	}
	if len(df.Body.Events) == 0 {
		return false
	}
	for _, event := range df.Body.Events {
		if event.EventName == "" {
			return false
		}
		if event.Params == "" {
			return false
		}
	}

	return true
}

func (df *DataFinder) AddUser(userUniqueId string) {
	if !df.isInit() {
		log.Warn("NOT INIT DataFinder")
		return
	}
	df.Body.User = &UserStruct{UserUniqueId: userUniqueId}
}

func (df *DataFinder) AddHeader(header *HeaderStruct) {
	if !df.isInit() {
		log.Warn("NOT INIT DataFinder")
		return
	}
	df.Body.Header = header
}

/*
	单次上传events数控制在20条以内，超过50条会报413
*/
func (df *DataFinder) AddEvent(event *EventStruct) {
	if !df.isInit() {
		log.Warn("NOT INIT DataFinder")
		return
	}
	event.LocalTimeMs = time.Now().UnixNano() / int64(time.Millisecond)
	df.Body.Events = append(df.Body.Events, event)
}

/*
	发起上报请求
*/
func (df *DataFinder) Request() (*ResponseBody, error) {
	if !df.valid() {
		return nil, RequestParamsValidErr
	}

	bodyByte, err := json.Marshal(df.Body)
	if err != nil {
		log.Error("marshal dataFinder body err", err)
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, DataRangersSinglePostUrl, strings.NewReader(string(bodyByte)))
	if err != nil {
		log.Error("create dataFinder http request err", err)
		return nil, err
	}

	// 添加header 标识
	request.Header.Add("Content-type", "application/json")
	request.Header.Add("X-MCS-AppKey", df.AppKey)

	resp, err := DefaultNetClient.Do(request)
	if err != nil {
		log.Error(request.URL, "send request err ", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != ResponseCodeOk {
		switch resp.StatusCode {
		case ResponseCodeMalformed:
			return nil, RequestMalformedErr
		case ResponseCodeTooMany:
			return nil, RequestTooManyErr
		case ResponseCodeFrequency:
			return nil, RequestFrequencyErr
		case ResponseCodeServerErr:
			return nil, RequestServerErr
		default:
			return nil, RequestServerErr
		}
	}
	respDataByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(request.URL, "read response data err", err)
		return nil, err
	}
	var responseBody ResponseBody
	err = json.Unmarshal(respDataByte, &responseBody)
	if err != nil {
		log.Error("marshal dataFinder response body err", err, string(respDataByte))
		return nil, err
	}
	return &responseBody, nil
}
