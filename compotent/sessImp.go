//session管理
package compotent

//session存储方式接口
type Session interface {
	Set(key string, value interface{}) error
	Get(key string) interface{}
	Delete(key string) error
	SessionID() string
}

//Session操作接口
type Provider interface {
	//初始化一个session, sid根据需要生成后传入
	SessionInit(sid string) (Session, error)
	//根据sid, 获取session
	SessionRead(sid string) (Session, error)
	//销毁session
	SessionDestroy(sid string) error
}
