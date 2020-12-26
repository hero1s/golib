package sdk

import (
	"encoding/json"
	"fmt"
	"time"
)

// generateData 生成data
func generateData(v interface{}, des3key string) (string, error) {
	originData, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return Encrypt(originData, des3key)
}

// generateMess 生成随机数
func generateMess() string {
	return fmt.Sprint(time.Now().Nanosecond())
}

// generateTimestamp 生成时间戳
func generateTimestamp() string {
	return fmt.Sprint(time.Now().Unix())
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	return fmt.Sprint(time.Now().UnixNano())
}

// BuildParams 封装请求信息
func BuildParams(v interface{}, appkey, des3key string) (map[string]string, error) {
	data, err := generateData(v, des3key)
	if err != nil {
		return nil, err
	}
	mess := generateMess()
	timestamp := generateTimestamp()
	sign := Sign(data, mess, timestamp, appkey)
	return map[string]string{
		"data":      data,
		"mess":      mess,
		"timestamp": timestamp,
		"sign":      sign,
		"sign_type": "sha256",
	}, nil
}

// BuildHeader 封装请求头
func BuildHeader(dealerID string) map[string]string {
	return map[string]string{
		"dealer-id":  dealerID,
		"request-id": generateRequestID(),
	}
}

// VerifySignature 验证签名是否一致
func VerifySignature(data, mess, timestamp, sign, appkey string) bool {
	originSign := Sign(data, mess, timestamp, appkey)
	return originSign == sign
}
