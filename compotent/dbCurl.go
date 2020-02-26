package compotent

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type CurdHandler struct {
	*gorm.DB
}

func (c *CurdHandler) FindItemByPk(out interface{}, id int) (err error) {
	c.First(out, id)
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("[query] query fail, panic:`%v`", p)
		}
	}()
	return
}

func (c *CurdHandler) FindItemByUk(out interface{}, where ...interface{}) (err error) {
	c.Find(out, where...)
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("[query] query fail, panic:`%v`", p)
		}
	}()
	return
}
