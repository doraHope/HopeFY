package settting

import (
	"github.com/doraHope/HopeFY/util/session"
	"log"

	"gopkg.in/ini.v1"

	_ "github.com/doraHope/HopeFY/middleware"
)

type App struct {
	JwSecret string //jwt密钥
}

type Server struct {
	Host string	//web服务域名
	Ip   string //web服务主机名
	Port string //web服务端口号
}

type Redis struct {
	Host		string
	Password 	string
	DB			int
}

type Database struct {
	Type        string //数据库类型
	User        string //用户
	Pw          string //密码
	Host        string //数据库主机名
	Port        string //端口
	Append      string //追加参数
	DB          string //数据库
	TablePrefix string //表前缀
}

type SessionManger struct {
	ServerName string
	CookieName string
	MaxLifeTime int64
}

var AppSetting = &App{}
var ServiceSetting = &Server{}
var DBSetting = &Database{}
var RedisSetting = &Redis{}
var SessionSetting = &SessionManger{}
var config *ini.File

func mapTo(section string, v interface{}) {
	err := config.Section(section).MapTo(v)
	if err != nil {
		log.Fatal("Cfg.MapTo RedisSetting err: %v", err)
	}
}

func Setup() {
	var err error
	config, err = ini.Load("../config/app.ini")
	if err != nil {
		log.Fatal("Fail to parse 'config/app.ini': %v", err)
	}
	mapTo("app", AppSetting)
	mapTo("database", DBSetting)
	mapTo("server", ServiceSetting)
	mapTo("redis", RedisSetting)
	mapTo("session", SessionSetting)
}


func RegisterAppMiddleware() {
	manager, err := session.NewSessionManager(SessionSetting.ServerName, SessionSetting.CookieName, SessionSetting.MaxLifeTime)
	if err != nil {
		//todo log
		log.Fatal(err)
	}
	session.Register(SessionSetting.ServerName, manager)
}