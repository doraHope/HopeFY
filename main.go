package main

import (
	"github.com/doraHope/HopeFY/router"
)

func main() {
	r := router.InitRouter()
	r.Run(":9090")
}
