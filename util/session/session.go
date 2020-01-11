package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/doraHope/HopeFY/components"
	"github.com/doraHope/HopeFY/enum"
	"github.com/doraHope/HopeFY/settting"
)

//session存储方式接口
type Session interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
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
	//回收
	SessionGC(maxLifeTime int64)
}

type Manager struct {
	cookieName       string
	lock             sync.Mutex //互斥锁
	provider         Provider   //存储session方式
	maxLifeTime      int64      //有效期
	sessionNumber    int64      //当前管理的会话数量
	maxSessionNumber int64      //最大会话数量
}

var provides = make(map[string]Provider) //并发web服务?

//实例化一个session管理器
func NewSessionManager(provideName, cookieName string, maxLifeTime int64) (*Manager, error) {
	provide, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q ", provideName)
	}
	return &Manager{cookieName: cookieName, provider: provide, maxLifeTime: maxLifeTime}, nil
}

//注册 由实现Provider接口的结构体调用
//这里为什么要provider作为传参, 而不是生成一个信息, 值得学习的编码方法
func Register(name string, provide Provider) {
	if provide == nil {
		panic("session: Register provide is nil")
	}
	if _, ok := provides[name]; ok {
		panic("session: Register called twice for provide " + name)
	}
	provides[name] = provide
}

//生成sessionId
func (manager *Manager) sessionId() string {
	b := make([]byte, 32)
	tryUp := 3
	tryNumber := 0
	for {
		if _, err := io.ReadFull(rand.Reader, b); err != nil {
			//todo log
			return ""
		}
		var session Session
		var err error
		if session, err = manager.provider.SessionRead(string(b)); err != nil {
			//todo log
			return ""
		}
		//如果出现重复sessionID, 则尝试生成新的sessionID
		if session != nil {
			//todo log
			tryNumber++
			if tryNumber >= tryUp {
				//todo log
				return ""
			}
			break
		}
		break
	}
	//加密
	return base64.URLEncoding.EncodeToString(b)
}

//判断当前请求的cookie中是否存在有效的session，存在返回，否则创建
func (manager *Manager) SessionStart(gc *gin.Context) (session Session) {
	manager.lock.Lock() //加锁
	defer manager.lock.Unlock()
	cookie, err := gc.Cookie(manager.cookieName)
	if err != nil || cookie == "" {
		//创建一个
		sid := manager.sessionId()
		if sid == "" {
			components.ResponseServError(gc, enum.SERVICE_ERROR, "")
		}
		session, _ = manager.provider.SessionInit(sid)
		gc.SetCookie(manager.cookieName, url.QueryEscape(sid), int(manager.maxLifeTime), "/", settting.ServiceSetting.Host, false, true)
		//当会话连接数量 > 连接上限则请求 启动一个协程清理
		manager.sessionNumber++
		if manager.sessionNumber >= manager.maxSessionNumber {
			go func() {
				manager.GC()
			}()
		}
	} else {
		sid, _ := url.QueryUnescape(cookie) //反转义特殊符号
		session, _ = manager.provider.SessionRead(sid)
	}
	return session
}

//销毁session 同时删除cookie
func (manager *Manager) SessionDestroy(gc *gin.Context) {
	cookie, err := gc.Cookie(manager.cookieName)
	if err != nil || cookie == "" {
		return
	} else {
		manager.lock.Lock()
		defer manager.lock.Unlock()
		sid, _ := url.QueryUnescape(cookie)
		manager.provider.SessionDestroy(sid)
		gc.SetCookie(manager.cookieName, url.QueryEscape(sid), -1, "/", settting.ServiceSetting.Host, false, true)
	}
}

//清理过期session
func (manager *Manager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionGC(manager.maxLifeTime)
	time.AfterFunc(time.Duration(manager.maxLifeTime), func() { manager.GC() })
}
