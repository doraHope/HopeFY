package bean

type Redis struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password omitempty"`
}
