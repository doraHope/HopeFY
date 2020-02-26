package middlerware

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"../config"
)

func init() {
	//root:chuyinwl01@(localhost:3306)/go?charset=utf8mb4&loc=Local
	//依赖配置项
	dsn := fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?%s",
		config.AppConfig.DB.User,
		config.AppConfig.DB.Password,
		config.AppConfig.DB.Host,
		config.AppConfig.DB.Port,
		config.AppConfig.DB.DBName,
		config.AppConfig.DB.Option,
	)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(fmt.Sprintf("数据库初始化异常, error:`%v`", err))
	}
	Cont.Register("db", db)
}

func GetDB(db interface{}) *gorm.DB {
	gd, ok := db.(*gorm.DB)
	if !ok {
		return nil
	}
	return gd
}
