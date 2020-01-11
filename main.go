package main

import (
	"github.com/doraHope/HopeFY/router"
	"github.com/doraHope/HopeFY/settting"
)

func main() {
	//注册中间件
	settting.RegisterAppMiddleware()
	//启动路由
	r := router.InitRouter()
	r.Run(":9090")
}
