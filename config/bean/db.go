package bean

type DB struct {
	Host     string `json:"host"`     //主机名
	Port     string `json:"port"`     //端口号
	User     string `json:"user"`     //用户名
	Password string `json:"password"` //密码
	DBName   string `json:"name"`     //数据库名
	Option   string `json:"opt"`      //配置项
	MaxOpen  int    `json:"max_open, omitempty"`
	MaxIdle  int    `json:"max_idle, omitempty"`
}
