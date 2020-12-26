package nacos

import (
	"fmt"
	"testing"
)

func TestAcm(t *testing.T) {
	// 从控制台命名空间管理的"命名空间详情"中拷贝 End Point、命名空间 ID
	var endpoint = "acm.aliyun.com"
	var namespaceId = "4d2e32d3-ffa4-4668-b213-495eecfaf27c"

	// 推荐使用 RAM 用户的 accessKey、secretKey
	var accessKey = "nafgL5OLBZiDgbfU"
	var secretKey = "zpfC1OBProcFxEhXbe7pEuMD5L0MO0"

	InitConfigClient(endpoint,namespaceId,accessKey,secretKey)

	PublishConfig("data_test1","group_test1","this is test1")

	GetConfig("data_test1","group_test1")

	ListenConfig("data_test1","group_test1", func(namespace, group, dataId, data string) {
		fmt.Println("ListenConfig group:" + group + ", dataId:" + dataId + ", data:" + data)
	})

	DeleteConfig("data_test1","group_test1")

}
