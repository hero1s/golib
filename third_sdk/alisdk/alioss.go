package alisdk

import (
	"errors"
	"fmt"
	"github.com/hero1s/golib/log"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/denverdino/aliyungo/ram"
	"github.com/denverdino/aliyungo/sts"
	"io"
	"path/filepath"
)

var (
	ossUrlRoot          = ""
	bucketName          = ""
	roleAcs             = ""
	accessKeyId_oss     = ""
	accessKeySecret_oss = ""
	region              = ""
)

var policyDocument = ram.PolicyDocument{
	Statement: []ram.PolicyItem{
		ram.PolicyItem{
			Action:   "oss:*",
			Effect:   "Allow",
			Resource: "acs:oss:*:*:*/anyprefix",
		},
	},
	Version: "1",
}

func InitAliOSS(accessKey, accessKeySecret, ossRoot, bucket, roleacs, ossRegion string) {
	accessKeyId_oss = accessKey
	accessKeySecret_oss = accessKeySecret
	ossUrlRoot = ossRoot
	bucketName = bucket
	roleAcs = roleacs
	region = ossRegion
}
func GetOssEndpoint() string {
	return ossUrlRoot
}
func GetOssPathUrl() string {
	return bucketName + "." + ossUrlRoot
}
func GetOssBucketName() string {
	return bucketName
}
func GetOssRegion() string {
	return region
}

// prePath: 存储目录结构的前缀目录，如 shop/title_image/
// id: 用id和prePath目录拼接起来做为目录返回
func GetImagePath(prePath, filename string) string {
	return fmt.Sprintf(filepath.Join(prePath, filename))
}

// 获取临时token授权(后续增加限制对应目录的权限 toney)
func GenerateOssToken(path string) (sts.AssumedRoleUserCredentials, error) {
	client := sts.NewClient(accessKeyId_oss, accessKeySecret_oss)
	var req sts.AssumeRoleRequest
	req.DurationSeconds = 1200 //20分钟
	req.RoleArn = roleAcs
	req.Policy = createOssPolicy(path)
	req.RoleSessionName = "client"
	log.Debugf("获取临时sts权限:%v", req.Policy)
	resp, err := client.AssumeRole(req)
	return resp.Credentials, err
}

// 权限控制策略文件
func createOssPolicy(path string) string {
	var policy string
	//policy = `{"Version":"1","Statement":[{"Effect":"Allow","Action":["oss:PutObject"],"Resource":["acs:oss:*:*:ram-test-app/usr001/*"]}]}`
	//policy = `{"Version":"1","Statement":[{"Effect":"Allow","Action":["oss:*"],"Resource":["acs:oss:*:*:*/*"]}]}`
	policy = `{"Version":"1","Statement":[{"Effect":"Allow","Action":["oss:*"],"Resource":["acs:oss:*:*:` + bucketName + "/" + path + `/*"]}]}`
	return policy
}

func PutFileToOSS(filename, newfilename string) error {
	client, err := oss.New(ossUrlRoot, accessKeyId_oss, accessKeySecret_oss)
	if err != nil {
		return errors.New("oss 创建客户端失败" + fmt.Sprintf("%v", err))
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return errors.New("oss 创建bucket失败" + fmt.Sprintf("%v", err))
	}
	err = bucket.PutObjectFromFile(newfilename, filename)
	if err != nil {
		return errors.New("oss 上传文件失败" + fmt.Sprintf("%v", err))
	}
	return nil
}
func PutFileStreamToOss(objectKey string, reader io.Reader) error {
	client, err := oss.New(ossUrlRoot, accessKeyId_oss, accessKeySecret_oss)
	if err != nil {
		return errors.New("OSS创建客户端失败" + fmt.Sprintf("%v", err))
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return errors.New("OSS创建bucket失败" + fmt.Sprintf("%v", err))
	}
	err = bucket.PutObject(objectKey, reader)
	if err != nil {
		return errors.New("OSS上传文件失败" + fmt.Sprintf("%v", err))
	}
	return nil
}
func DelFileFromOSS(filename string) error {
	client, err := oss.New(ossUrlRoot, accessKeyId_oss, accessKeySecret_oss)
	if err != nil {
		return errors.New("oss 创建客户端失败" + fmt.Sprintf("%v", err))
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return errors.New("oss 创建bucket失败" + fmt.Sprintf("%v", err))
	}
	err = bucket.DeleteObject(filename)
	if err != nil {
		return errors.New("oss 删除object失败" + fmt.Sprintf("%v", err))
	}
	return nil
}
