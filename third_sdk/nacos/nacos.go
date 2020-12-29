package nacos

import (
	"encoding/json"
	"github.com/hero1s/golib/log"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var iClient config_client.IConfigClient = nil

func InitConfigClient(endpoint, namespaceId, accessKey, secretKey string) error {
	clientConfig := constant.ClientConfig{
		Endpoint:       endpoint + ":8080",
		NamespaceId:    namespaceId,
		AccessKey:      accessKey,
		SecretKey:      secretKey,
		TimeoutMs:      5 * 1000,
		ListenInterval: 30 * 1000,
	}
	var err error
	iClient, err = clients.CreateConfigClient(map[string]interface{}{
		"clientConfig": clientConfig,
	})
	if err != nil {
		log.Errorf("创建nacos客户端错误:%v", err)
		return err
	}
	return nil
}

//发布配置
func PublishConfig(dataId, group, content string) error {
	success, err := iClient.PublishConfig(vo.ConfigParam{
		DataId:  dataId,
		Group:   group,
		Content: content})

	if success {
		log.Info("Publish config successfully.")
	} else {
		log.Errorf("Publish config error:%v.", err)
	}
	return err
}

//获取配置
func GetConfig(dataId, group string) (string, error) {
	content, err := iClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group})
	log.Infof("Get config：%v,err:%v", content, err)
	return content, err
}

func GetConfigToStruct(dataId, group string, cfg interface{}) error {
	content, err := iClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group})
	log.Infof("Get config：%v,err:%v", content, err)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(content), &cfg)
	if err != nil {
		log.Errorf("读取配置文件解析失败:%v,%v--%v", dataId, group, err)
	}
	return err
}

//监听配置
func ListenConfig(dataId, group string, f func(namespace, group, dataId, data string)) {
	iClient.ListenConfig(vo.ConfigParam{
		DataId:   dataId,
		Group:    group,
		OnChange: f,
	})
}

//删除配置
func DeleteConfig(dataId, group string) {
	success, err := iClient.DeleteConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group})

	if success {
		log.Info("Delete config successfully.")
	} else {
		log.Errorf("Delete config fail:%v.", err)
	}
}
