package proxy

import (
	"encoding/json"
	"io/ioutil"
)

// 配置文件
type Config struct {
	ProxyHost  []Upstream `json:"proxy_host"`  //代理配置
	ServerPort string     `json:"server_port"` // 代理转发服务启动的端口
}

type Upstream struct {
	Upstream string `json:"upstream"`  //远端地址
	UpHost   string `json:"up_host"`   //代理地址
	Path     string `json:"path"`      //转发路径
	TrimPath bool   `json:"trim_path"` //移除路径
	IsAuth   bool   `json:"is_auth"`   //是否鉴权
	AuthKey  string `json:"auth_key"`  //简单的校验Key
}

func LoadConfig(fileName string) (Config, error) {
	var cfg Config
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return cfg, err
	}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
