package response

// 统一 json 结构体
type JsonObject struct {
	Code string      `json:"code" desc:"状态码"`
	Data interface{} `json:"content,omitempty" desc:"数据"`
	Msg  string      `json:"msg" desc:"消息"`
}

// 全局分页对象
type PageBean struct {
	Page     int         `json:"page" desc:"当前页"`
	PageSize int         `json:"page_size" desc:"每页最大行数"`
	Total    int         `json:"total" desc:"总记录数"`
	Rows     interface{} `json:"rows" desc:"每行数据"`
}
