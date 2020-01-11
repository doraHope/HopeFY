package middleware

import (
    "encoding/json"
    "fmt"
    "github.com/doraHope/HopeFY/util/common"
	"github.com/doraHope/HopeFY/util/session"
	"log"
	"sync"
	"time"

    "github.com/go-redis/redis"

    "github.com/doraHope/HopeFY/settting"
)

type RedisProvider struct {
	Client  *redis.Client
	lock    sync.Mutex
    stNumber map[string]int64
}

//session实现
type SessionStore struct {
	sid              string                         `json:"sid"`              //会话id
	LastAccessedTime time.Time                      `json:"lastAccessedTime"` //最后访问时间
	value            map[interface{}]interface{} `json:"data, omitempty"`  //session 里面存储的值
}

var redisProvider *RedisProvider

func init() {
	redisProvider = &RedisProvider{
		Client: redis.NewClient(&redis.Options{
			Addr:     settting.RedisSetting.Host,
			Password: settting.RedisSetting.Password,
			DB:       settting.RedisSetting.DB,
		}),
	}
	pong, err := redisProvider.Client.Ping().Result()
	if err != nil {
		//todo log
		log.Fatal(fmt.Sprintf("connect redis error, %v\n", err))
	}
	fmt.Printf("ping redis success, %v\n", pong)
	//todo log info
}

//设置
func (st *SessionStore) Set(key, value interface{}) error {
	st.value[key] = value
	redisProvider.SessionUpdate(st.sid)
	return nil
}

//获取session
func (st *SessionStore) Get(key interface{}) interface{} {
	redisProvider.SessionUpdate(st.sid)
	if v, ok := st.value[key]; ok {
		return v
	} else {
		return nil
	}
	return nil
}

//删除
func (st *SessionStore) Delete(key interface{}) error {
	delete(st.value, key)
	redisProvider.SessionUpdate(st.sid)
	return nil
}

//获取sessionID
func (st *SessionStore) SessionID() string {
	return st.sid
}

func (rp *RedisProvider) SessionInit(sid string) (session.Session, error) {
	rp.lock.Lock()
	defer rp.lock.Unlock()
	v := make(map[interface{}]interface{}, 0)
	timestamp := time.Now()
	session := &SessionStore{
	    sid: sid,
	    LastAccessedTime: timestamp,
	    value: v,
	}
	rp.stNumber[sid] = timestamp.Unix()
	return session, nil
}

func (rp *RedisProvider) parseSessionStore(sid string) (session.Session, error) {
    st := &SessionStore{}
    var err error
    var element string
    if element, err = rp.Client.Get(sid).Result(); err == nil {
        err = json.Unmarshal([]byte(element), st)
        if err == nil {
            //todo log
            return st, nil
        }
    }
    return nil, fmt.Errorf("[session] 解析失败, %v", err)
}

func (rp *RedisProvider) SessionRead(sid string) (session.Session, error) {
	return rp.parseSessionStore(sid)
}

func (rp *RedisProvider) SessionDestroy(sid string) error {
    _, err := rp.Client.Del(sid).Result()
    return err
}

func (rp *RedisProvider) SessionGC(maxLifeTime int64) {
	rp.lock.Lock()
	defer rp.lock.Unlock()
	//对话会话进行排序
    result := common.SortMapSI64(rp.stNumber)
    //遍历排序后的结果, 并清除对应的会话
	for _, sid := range result {
	    if _, ok := rp.stNumber[sid]; ok {
            session, err := rp.parseSessionStore(sid)
            if err != nil {
                //todo log
                continue
            }
            st, ok := session.(*SessionStore)
            if !ok {
                //todo log
                continue
            }
            //经过排序清楚过期会话
            if st.LastAccessedTime.Unix() < time.Now().Unix() {
                st.Delete(sid)
            } else {
                break;
            }
        }

    }
}
func (rp *RedisProvider) SessionUpdate(sid string) error {
    session, err := rp.parseSessionStore(sid)
    if err != nil {
        //todo log
        return err
    }
    st, ok := session.(*SessionStore)
    if ok {
        //todo log
        return fmt.Errorf("[session] SessionStore类型断言异常")
    }
    st.LastAccessedTime = time.Now()
    jsonSt, err := json.Marshal(st)
    if err != nil {
        return fmt.Errorf("[session] SessionStore转换为json失败, %v", err)
    }
    _, err = rp.Client.Set(sid, jsonSt, 2 * time.Hour).Result();
	return fmt.Errorf("[session] SessionStore缓存到redis失败, %v", err)
}
