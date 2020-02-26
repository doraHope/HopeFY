//数据CURD, 在考虑要不要这个东西
package compotent

import (
	_ "../middlerware"
)

type Curd interface {
	FindItemByPk(interface{}, int) error
	FindItemByUk(interface{}, string, ...interface{}) error
	updateByPk(int) (int, error)
	updateByUk(string, ...interface{}) (int, error)
	save(interface{}) (int, error)
}
