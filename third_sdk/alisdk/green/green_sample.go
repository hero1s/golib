package green_sdk

import (
	"encoding/json"
	"git.moumentei.com/plat_go/golib/log"
	greensdksample2 "git.moumentei.com/plat_go/golib/third_sdk/alisdk/green/greensdksample"
	uuid2 "git.moumentei.com/plat_go/golib/third_sdk/alisdk/green/uuid"
	"path"
	"regexp"
)

var accessKeyId = "<your access key id>"
var accessKeySecret = "<your access key secret>"

func InitGreenSdk(KeyId, KeySecret string) {
	accessKeyId = KeyId
	accessKeySecret = KeySecret
}

type Detail struct {
	Label      string  `json:"label"`
	Rate       float64 `json:"rate"`
	Scene      string  `json:"scene"`
	Suggestion string  `json:"suggestion"`
}

type ScanDataResp struct {
	Code    int64    `json:"code"`
	Content string   `json:"content"`
	Url     string   `json:"url"`
	DataId  string   `json:"dataId"`
	Msg     string   `json:"msg"`
	TaskId  string   `json:"taskId"`
	Results []Detail `json:"results"`
}

type ScanResp struct {
	Code      int64          `json:"code"`
	Msg       string         `json:"msg"`
	RequestId string         `json:"requestId"`
	Data      []ScanDataResp `json:"data"`
}

//鉴定图片

/*porn：图片智能鉴黄
terrorism：图片暴恐涉政识别
ad：图文违规识别
qrcode：图片二维码识别
live：图片不良场景识别
logo：图片logo识别
[]string{"porn","terrorism","ad","live","qrcode","logo"}
*/

func CheckImageScan(imageUrl string, scenes []string) bool {
	profile := greensdksample2.Profile{AccessKeyId: accessKeyId, AccessKeySecret: accessKeySecret}

	path := "/green/image/scan"
	//是否视频
	if CheckFileNameIsVideo(imageUrl) { //视频文件后台自动冻结
		return true
	}
	clientInfo := greensdksample2.ClinetInfo{Ip: "127.0.0.1"}
	// 构造请求数据
	bizType := "Green"
	task := greensdksample2.Task{DataId: uuid2.Rand().Hex(), Url: imageUrl}
	tasks := []greensdksample2.Task{task}

	bizData := greensdksample2.BizData{bizType, scenes, tasks}

	var client greensdksample2.IAliYunClient = greensdksample2.DefaultClient{Profile: profile}

	// your biz code
	strResp := client.GetResponse(path, clientInfo, bizData)
	log.Debugf("色情图片检测:%v", strResp)
	var resp ScanResp
	if err := json.Unmarshal([]byte(strResp), &resp); err != nil {
		log.Errorf("解析鉴黄结构体错误:%v", err)
	} else {
		log.Debugf("鉴黄结果返回:%+v", resp)
		for _, d := range resp.Data {
			for _, res := range d.Results {
				if res.Suggestion == "block" {
					log.Errorf("鉴定为屏蔽:%v,场景:%+v", d.Url, res.Scene)
					return false
				}
				if res.Rate > 80 && res.Label == "porn" {
					log.Errorf("鉴定场景分数:%+v", res)
					return false
				}
			}
		}
	}
	return true

}

//鉴定文本
func CheckTextScan(text string) bool {
	profile := greensdksample2.Profile{AccessKeyId: accessKeyId, AccessKeySecret: accessKeySecret}

	pattern := "^[A-Za-z0-9]+$"
	ok, _ := regexp.MatchString(pattern, text)
	if ok {
		log.Debugf("纯字母数字文本不检测:%v", text)
		return true
	}
	if len(text) < 4 {
		return true
	}

	path := "/green/text/scan"
	clientInfo := greensdksample2.ClinetInfo{Ip: "127.0.0.1"}

	// 构造请求数据
	bizType := "Green"
	scenes := []string{"antispam"}

	task := greensdksample2.Task{DataId: uuid2.Rand().Hex(), Content: text}
	tasks := []greensdksample2.Task{task}

	bizData := greensdksample2.BizData{bizType, scenes, tasks}

	var client greensdksample2.IAliYunClient = greensdksample2.DefaultClient{Profile: profile}

	// your biz code
	strResp := client.GetResponse(path, clientInfo, bizData)
	log.Debugf("色情文字检测:%v", strResp)
	var resp ScanResp
	if err := json.Unmarshal([]byte(strResp), &resp); err != nil {
		log.Errorf("解析鉴黄结构体错误:%v", err)
	} else {
		log.Debugf("鉴黄结果返回:%+v", resp)
		for _, d := range resp.Data {
			for _, res := range d.Results {
				if res.Suggestion == "block" {
					log.Errorf("鉴定为屏蔽:%v", d.Content)
					return false
				}
			}
		}
	}
	return true
}

//判断文件是否视频
func CheckFileNameIsVideo(fileName string) bool {
	filenameWithSuffix := path.Base(fileName)  //获取文件名带后缀
	fileSuffix := path.Ext(filenameWithSuffix) //获取文件后缀
	log.Debugf("文件后缀名:%v", fileSuffix)
	if fileSuffix == ".mp4" {
		return true
	}

	return false
}
