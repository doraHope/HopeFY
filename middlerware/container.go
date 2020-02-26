//容器, 用于管理中间件
package middlerware

import (
	"fmt"
	"sync"
)

type Container struct {
	mux sync.Mutex
	mds map[string]interface{}
}

var Cont *Container = &Container{
	sync.Mutex{},
	make(map[string]interface{}, 8),
}

func (c *Container) Register(n string, sl interface{}) {
	c.mux.Lock()
	defer c.mux.Unlock()
	_, ok := c.mds[n]
	if !ok {
		c.mds[n] = sl
	}
}

func (c *Container) Get(n string) (interface{}, error) {
	sl, ok := c.mds[n]
	if !ok {
		return nil, fmt.Errorf("未找到单例`%v`", n)
	}
	return sl, nil
}
