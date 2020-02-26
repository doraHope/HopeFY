package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	b "./bean"
)

type App struct {
	DB      b.DB      `json:"mysql"`
	Session b.Session `json:"session"`
	Redis   b.Redis   `json:"redis"`
	Logger  b.Logger  `json:"logger"`
}

var AppConfig App

var cfn = flag.String("c", `D:\go\stayWithYou\swy\app.json`, "应用程序配置文件")

func init() {
	flag.Parse()
	fd, err := os.OpenFile(*cfn, os.O_RDONLY, 644)
	if err != nil {
		panic(fmt.Sprintf("文件打`%v`开失败, err:`%v`\n", *cfn, err))
	}
	cf, err := ioutil.ReadAll(fd)
	if err != nil {
		panic(fmt.Sprintf("配置读取异常, err:`%v`\n", err))
	}
	loadAppConfig(cf)
}

func loadAppConfig(config []byte) error {
	err := json.Unmarshal(config, &AppConfig)
	if err != nil {
		return fmt.Errorf("初始化应用配置失败, error:`%v`", err)
	}
	return nil
}
