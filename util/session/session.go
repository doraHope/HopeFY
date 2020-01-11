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
)

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
	//回收
	SessionGC(maxLifeTime int64) int
}

type Manager struct {
	cookieName       string
	domain           string
	lock             sync.Mutex //互斥锁
	provider         Provider   //存储session方式
	maxLifeTime      int64      //有效期
	sessionNumber    int        //当前管理的会话数量
	maxSessionNumber int        //最大会话数量
}

var provides = make(map[string]Provider) //并发web服务?

//实例化一个session管理器
func NewSessionManager(provideName, cookieName, domain string, maxLifeTime int64, maxSessionNumber int) (*Manager, error) {
	provide, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q ", provideName)
	}
	manager := &Manager{
		cookieName:       cookieName,
		domain:           domain,
		provider:         provide,
		maxLifeTime:      maxLifeTime,
		maxSessionNumber: maxSessionNumber,
	}
	time.AfterFunc(time.Duration(maxLifeTime), func() {
		manager.GC()
	})
	return manager, nil
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
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		//todo log
		return ""
	}

	//加密
	return base64.URLEncoding.EncodeToString(b)
}

//设置会话标识
func (manager *Manager) setSession() (Session, string, error) {
	//创建一个
	var sid string
	tryNumber := 0
	tryUp := 3
	for {
		sid = manager.sessionId()
		if sid == "" {
			return nil, "", fmt.Errorf("[session] 产生空的sessionID")
		}
		var session Session
		var err error
		if session, err = manager.provider.SessionRead(string(sid)); err != nil {
			//todo log
			return nil, "", fmt.Errorf("[session] 取出session出错, %v", err)
		}
		//如果出现重复sessionID, 则尝试生成新的sessionID
		if session != nil {
			//todo log
			tryNumber++
			if tryNumber >= tryUp {
				//todo log
				return nil, "", fmt.Errorf("[session] 3次连续取出重复的sessionID")
			}
			continue
		}
		break
	}
	session, err := manager.provider.SessionInit(sid)
	if err != nil {
		return nil, "", fmt.Errorf("[session] session初始化失败, %v", err)
	}
	return session, sid, nil
}

//判断当前请求的cookie中是否存在有效的session，存在返回，否则创建
func (manager *Manager) SessionStart(gc *gin.Context) (Session, error) {
	manager.lock.Lock() //加锁
	defer manager.lock.Unlock()
	cookie, err := gc.Cookie(manager.cookieName)
	if err != nil || cookie == "" {
		err = err
		session, sid, err := manager.setSession()
		if err != nil {
			return nil, fmt.Errorf("[manager] 设置会话标识失败")
		}
		//当会话连接数量 > 连接上限则请求 启动一个协程清理
		manager.sessionNumber++
		if manager.sessionNumber >= manager.maxSessionNumber {
			go func() {
				manager.sessionNumber = manager.provider.SessionGC(manager.maxLifeTime)
			}()
		}
		gc.SetCookie(manager.cookieName, url.QueryEscape(sid), int(manager.maxLifeTime), "/", manager.domain, false, true)
		return session, nil
	}
	sid, _ := url.QueryUnescape(cookie) //反转义特殊符号
	session, err := manager.provider.SessionRead(sid)
	if err != nil {
		return nil, fmt.Errorf("[manager] session读取失败, %v", err)
	}
	if session == nil {
		session, sid, err := manager.setSession()
		if err != nil {
			return nil, fmt.Errorf("[manager] 设置会话标识失败")
		}
		//当会话连接数量 > 连接上限则请求 启动一个协程清理
		manager.sessionNumber++
		if manager.sessionNumber >= manager.maxSessionNumber {
			go func() {
				manager.sessionNumber = manager.provider.SessionGC(manager.maxLifeTime)
			}()
		}
		gc.SetCookie(manager.cookieName, url.QueryEscape(sid), int(manager.maxLifeTime), "/", manager.domain, false, true)
		return session, nil
	}
	return session, nil
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
		gc.SetCookie(manager.cookieName, url.QueryEscape(sid), -1, "/", manager.domain, false, true)
	}
}

//清理过期session
func (manager *Manager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.sessionNumber = manager.provider.SessionGC(manager.maxLifeTime)
	time.AfterFunc(time.Duration(manager.maxLifeTime), func() {
		manager.GC()
	})
}
