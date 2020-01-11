package router

import (
	"github.com/doraHope/HopeFY/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() *gin.Engine {
	router := gin.New()
	apiUser := router.Group("/user")
	{
		apiUser.GET("/:id?action=login", user.Action)
        apiUser.GET("", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{
                "code": 0,
                "msg":  "success",
                "data": nil,
            })
        })
	}
	return router
}
