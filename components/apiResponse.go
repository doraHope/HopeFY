package components

import (
	"github.com/doraHope/HopeFY/enum"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ResponseOk(gc *gin.Context, code int, message string, data map[string]interface{}) {
	if msg := enum.ErrMsg(code); message != "" {
		message = msg
	}
	gc.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  message,
		"data": data,
	})
}

func ResponseServError(gc *gin.Context, code int, message string) {
	if msg := enum.ErrMsg(code); message != "" {
		message = msg
	}
	gc.JSON(http.StatusInternalServerError, gin.H{
		"code": code,
		"msg":  message,
		"data": nil,
	})
}
