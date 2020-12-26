package alisdk

import (
	"encoding/json"
	"fmt"
	"github.com/hero1s/golib/helpers/encode"
	"github.com/hero1s/golib/log"
	"github.com/denverdino/aliyungo/common"
	. "github.com/denverdino/aliyungo/opensearch"
	"net/http"
	"time"
)

var (
	region_search          = ""
	accessKeyId_search     = ""
	accessKeySecret_search = ""
)
var (
	CMD_ADD    = "add"
	CMD_UPDATE = "update"
	CMD_DELTE  = "delete"
)

type SearchResult struct {
	Searchtime float64       `json:"searchtime"`
	Total      int64         `json:"total"`
	Num        int64         `json:"num"`
	Viewtotal  int64         `json:"viewtotal"`
	Items      []interface{} `json:"items"`
}
type SearchResp struct {
	Status    string       `json:"status"`
	RequestId string       `json:"request_id"`
	Result    SearchResult `json:"result"`
	Errors    []*Error     `json:"errors"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

//map[searchtime:0.030263 total:1 num:1 viewtotal:1 items:[map[index_name:jiangbei_product description:小商品 id:1 name:toney]]

func InitAliOpenSearch(region, accessKey, accessKeySecret string) {
	region_search = region
	accessKeyId_search = accessKey
	accessKeySecret_search = accessKeySecret
}
func newClient() *Client {
	client := NewClient(Internet, common.Region(region_search), accessKeyId_search, accessKeySecret_search)
	return client
}

// 获取应用状态
func GetSearchStatus(appName string) (interface{}, error) {
	client := newClient()
	var resp interface{}
	err := client.GetStatus(appName, &resp)
	log.Debugf(fmt.Sprintf("应用状态:%v--err:%v", resp, err))
	return resp, err
}

// 上传文档(支持新增,更新，删除的批量操作)
func PushSearch(appName, tableName string, items interface{}, cmd string) (interface{}, error) {
	client := newClient()
	args := PushArgs{}
	args.Table_name = tableName
	type AddItem struct {
		Cmd       string      `json:"cmd"`
		Timestamp int64       `json:"timestamp"`
		Fields    interface{} `json:"fields"`
	}
	var add [1]AddItem
	add[0].Timestamp = time.Now().Unix()
	add[0].Cmd = cmd
	add[0].Fields = items

	str, _ := json.Marshal(&add)
	args.Items = string(str)
	log.Debugf(args.Items)
	var resp interface{}
	err := client.Push(appName, args, &resp)
	log.Debugf(fmt.Sprintf("上传文档:%v--err:%v", resp, err))
	return resp, err
}
func PushSearchMore(appName, tableName string, cmd string, items ...interface{}) (interface{}, error) {
	client := newClient()
	args := PushArgs{}
	args.Table_name = tableName
	type AddItem struct {
		Cmd       string      `json:"cmd"`
		Timestamp int64       `json:"timestamp"`
		Fields    interface{} `json:"fields"`
	}
	var add []AddItem
	for _, v := range items {
		var a AddItem
		a.Timestamp = time.Now().Unix()
		a.Cmd = cmd
		a.Fields = v
		add = append(add, a)
	}
	str, _ := json.Marshal(&add)
	args.Items = string(str)
	log.Debugf(args.Items)
	var resp interface{}
	err := client.Push(appName, args, &resp)
	if err != nil {
		log.Infof(fmt.Sprintf("上传文档:%v,%v--err:%v", args, resp, err))
	}
	return resp, err
}

//搜索(非必选参数使用默认)
func Search(args SearchArgs) (SearchResp, error) {
	client := newClient()
	var resp interface{}
	var search SearchResp
	err := client.Search(args, &resp)

	if err == nil {
		err = encode.ChangeStructByEncodeJson(resp, &search)
	} else {
		log.Infof("search error:%v", err)
	}
	if search.Status != "OK" {
		log.Infof(fmt.Sprintf("search FAIL:%+v,%+v", args, search))
	}
	return search, err
}

//下拉提示
type Suggestion struct {
	Suggestion string `json:"suggestion"`
}

func Suggest(appName, suggestName, query string, filterSearch string) []string {
	client := newClient()
	var resp interface{}
	type SuggestArgs struct {
		//搜索主体
		Query string `ArgName:"query"`
		//要查询的应用名
		Index_name   string `ArgName:"index_name"`
		Suggest_name string `ArgName:"suggest_name"`
		Hit          string `ArgName:"hit"`
	}
	args := SuggestArgs{Query: query, Index_name: appName, Suggest_name: suggestName, Hit: "5"}
	client.InvokeByAnyMethod(http.MethodGet, "", "/suggest", args, &resp)
	type SuggestResp struct {
		Suggestions []Suggestion
	}
	var sug SuggestResp
	encode.ChangeStructByEncodeJson(resp, &sug)
	var res []string
	var flagFilter = false
	if len(filterSearch) > 4 {
		flagFilter = true
	}
	for _, v := range sug.Suggestions {
		if flagFilter == true {
			query := fmt.Sprintf("query=default:%s%s", v.Suggestion, filterSearch)
			resp1, err1 := Search(SearchArgs{Index_name: appName, Query: query})
			if err1 == nil && resp1.Result.Num > 0 {
				res = append(res, v.Suggestion)
			}
		} else {
			res = append(res, v.Suggestion)
		}
	}
	log.Debugf(fmt.Sprintf("suggest:%+v----》%+v---%v", resp, sug, res))
	return res
}
