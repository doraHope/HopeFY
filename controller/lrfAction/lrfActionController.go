package lrfAction

import (
	"github.com/doraHope/HopeFY/components"
	"github.com/doraHope/HopeFY/enum"
	"github.com/gin-gonic/gin"
)

type LoginForm struct {
	UserAccNo    string `form:"acc_no"`
	UserPassword string `form:"password"`
	VerifyCode   string `form:verify_code`
}

func Login(gc *gin.Context) {
	lf := &LoginForm{}
	err := gc.ShouldBind(lf)
	if err != nil {
		//todo log
		components.ResponseOk(gc, enum.INVALID_PARAMS, "", nil)
	}

}
