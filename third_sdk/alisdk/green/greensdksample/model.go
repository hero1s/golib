package greensdksample

type ClinetInfo struct {
	SdkVersion  string `json:"sdkVersion"`
	CfgVersion  string `json:"cfgVersion"`
	UserType    string `json:"userType"`
	UserId      string `json:"userId"`
	UserNick    string `json:"userNick"`
	Avatar      string `json:"avatar"`
	Imei        string `json:"imei"`
	Imsi        string `json:"imsi"`
	Umid        string `json:"umid"`
	Ip          string `json:"ip"`
	Os          string `json:"os"`
	Channel     string `json:"channel"`
	HostAppName string `json:"hostAppName"`
	HostPackage string `json:"hostPackage"`
	HostVersion string `json:"hostVersion"`
}

type Task struct {
	DataId  string `json:"dataId"`
	Url     string `json:"url"`
	Content string `json:"content"`
}

type BizData struct {
	BizType string   `json:"bizType"`
	Scenes  []string `json:"scenes"`
	Tasks   []Task   `json:"tasks"`
}
