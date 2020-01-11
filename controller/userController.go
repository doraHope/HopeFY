package user

import (
    "fmt"
    "github.com/doraHope/HopeFY/components"
    "github.com/doraHope/HopeFY/enum"
    "github.com/doraHope/HopeFY/settting"
    "github.com/gin-gonic/gin"
)

func SayHello(gc *gin.Context) {
    session, err := settting.SManager.SessionStart(gc)
    if err != nil {
        components.ResponseOk(gc, enum.SERVICE_ERROR, "", nil)
    } else {
        components.ResponseOk(gc, enum.SUCCESS, "success", nil)
    }
    //session.Set("dora", "hope for you")
    //session.Set("kuai", "hope for you")
    //session.Delete("kuai")
    //fmt.Println("hello world")
    _ = session
}


