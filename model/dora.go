package model

import (
	"fmt"

	"../compotent"
	"../middlerware"
)

type DoraModel struct {
	compotent.CurdHandler        //数据库连接实例
	tn                    string //数据库名
}

var DoraHandler DoraModel

func init() {
	db, err := middlerware.Cont.Get("db")
	if err != nil {
		panic(fmt.Sprintf("内部错误, error:`%v`", err))
	}
	gd := middlerware.GetDB(db) //形成依赖
	if gd == nil {
		panic("内部错误, 数据库单例类型异常")
	}
	tn := "dora"
	DoraHandler = DoraModel{
		compotent.CurdHandler{
			gd.Table(tn),
		},
		tn,
	}
}
