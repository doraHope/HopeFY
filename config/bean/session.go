package bean

type Session struct {
	Prefix         string `json:"prefix"`         //k-v缓存中key的前缀
	Expired        int    `json:"expired"`        //有效时间, 单位m
	MaxOpenSession int    `json:"maxOpenSession"` //最大会话数量
}
