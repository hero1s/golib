package qijin

import (
	"encoding/json"
	"fmt"
	"github.com/hero1s/golib/helpers/crypto"
	"github.com/hero1s/golib/helpers/encode"
	"github.com/hero1s/golib/helpers/http_client"
	"github.com/hero1s/golib/log"
	"net/http"
	"strconv"
	"time"
)

//才薪-财税系统sdk

/*  接口名称	 */
const (
	gateWayTest = "https://www.qipaiwork.com/api/test"
	gateWay     = "https://www.qipaiwork.com/api/api"

	//用户创建
	methodUrl_userCreation = "/qjapi/userCreation"
	//打款
	methodUrl_orderCreation = "/qjapi/paymentOrderCreation"
	//查询
	methodUrl_orderDetail = "/qjapi/paymentOrderDetail"
)

type QiJinClient struct {
	Mid                string `json:"mid" desc:"商户编号"`
	MerchantPrivateKey []byte `json:"merchantPrivateKey" desc:"商户私钥"`
	SystemPublicKey    string `json:"systemPublicKey" desc:"系统公钥"`
	GateWay            string `json:"gate_way" desc:"网关地址"`
}

type ApiCommonHead struct {
	MerchantNumber string
	Version        string
	Timestamp      string
	Sign           string
}

var defaultClient *QiJinClient

func GetClient() *QiJinClient {
	return defaultClient
}

func InitQiJin(mid string, privateKey []byte, online bool) {
	defaultClient = new(QiJinClient)
	defaultClient.MerchantPrivateKey = privateKey
	defaultClient.Mid = mid
	defaultClient.GateWay = gateWayTest
	if online {
		defaultClient.GateWay = gateWay
	}
}

func newApiCommonHead(c *QiJinClient) ApiCommonHead {
	head := ApiCommonHead{
		MerchantNumber: c.Mid,
		Version:        "V1.0.0",
		Timestamp:      strconv.FormatInt(time.Now().UnixNano()/1e6, 10),
		Sign:           "",
	}
	return head
}

func addHead(h ApiCommonHead) http.Header {
	headers := http.Header{}
	headers.Add("drg-access-sign", h.Sign)
	headers.Add("drg-access-timestamp", h.Timestamp)
	headers.Add("drg-access-version", h.Version)
	headers.Add("drg-access-merchant", h.MerchantNumber)
	headers.Add("Content-Type", "application/json;charset=utf-8")

	return headers
}

//创建用户
func (c *QiJinClient) UserCreation(user UserCreation) (UserCreationResponse, error) {
	body, err := c.postDataHttps(methodUrl_userCreation, encode.ChangeStructToJsonMap(user))
	var resp UserCreationResponse
	if err != nil {
		return resp, err
	}
	err = json.Unmarshal(body, &resp)
	return resp, err
}

//打款订单
func (c *QiJinClient) OrderCreation(order OrderCreation) (OrderCreationResponse, error) {
	body, err := c.postDataHttps(methodUrl_orderCreation, encode.ChangeStructToJsonMap(order))
	var resp OrderCreationResponse
	if err != nil {
		return resp, err
	}
	err = json.Unmarshal(body, &resp)
	return resp, err
}

//订单查询
func (c *QiJinClient) OrderDetail(order OrderDetail) (OrderDetailResponse, error) {
	body, err := c.postDataHttps(methodUrl_orderDetail, encode.ChangeStructToJsonMap(order))
	var resp OrderDetailResponse
	if err != nil {
		return resp, err
	}
	err = json.Unmarshal(body, &resp)
	return resp, err
}

func (c *QiJinClient) postDataHttps(url string, data map[string]interface{}) ([]byte, error) {
	url = c.GateWay + url
	commonHead := newApiCommonHead(c)

	parameterStr, _ := json.Marshal(data)
	digest, err := crypto.Md5WithRsa(fmt.Sprintf(`%v%v%v%v`,
		commonHead.Timestamp, commonHead.Version, commonHead.MerchantNumber, string(parameterStr)), c.MerchantPrivateKey)
	if err != nil {
		log.Errorf("RSA错误:%v", err)
		return nil, err
	}
	commonHead.Sign = digest
	rsp, err := http_client.HttpsRequest(http_client.Request{
		Method: "POST",
		URL:    url,
		Header: addHead(commonHead),
		Body:   string(parameterStr),
	})
	log.Debugf("resp:%v,err:%v", string(rsp), err)
	return rsp, err
}
