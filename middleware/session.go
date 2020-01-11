package middleware

import (
    "encoding/json"
    "fmt"
    "github.com/doraHope/HopeFY/util/common"
	"github.com/doraHope/HopeFY/util/session"
	"sync"
	"time"

    "github.com/go-redis/redis"
)

type RedisProvider struct {
	Client  *redis.Client
	lock    sync.Mutex
    stNumber map[string]int64
}

//session实现
type SessionStore struct {
	Sid              string                         `json:"sid"`              //会话id
	LastAccessedTime time.Time                      `json:"lastAccessedTime"` //最后访问时间
	Data             map[string]interface{} `json:"data, omitempty"`  //session 里面存储的值
}

var redisProvider *RedisProvider

func Setup(host, password string, db int) (*RedisProvider, error) {
	stN := make(map[string]int64, 64)
	redisProvider = &RedisProvider{
		Client: redis.NewClient(&redis.Options{
			Addr:     host,
			Password: password,
			DB:       db,
		}),
		stNumber: stN,
	}
	_, err := redisProvider.Client.Ping().Result()
	if err != nil {
		//todo log
		return nil, err
	}
	//todo log info
	return redisProvider, nil
}

func NewRedisProvider() () {

}

func init() {


}

func NewSessionStore(sid string, lastACT time.Time, data map[string]interface{}) *SessionStore {
	if data == nil {
		data = make(map[string]interface{}, 8)
	}
	return &SessionStore{
		Sid: sid,
		LastAccessedTime: lastACT,
		Data: data,
	}
}

//设置
func (st *SessionStore) Set(key string, value interface{}) error {
	st.Data[key] = value
	redisProvider.SessionUpdate(st)
	return nil
}

//获取session
func (st *SessionStore) Get(key string) interface{} {
	redisProvider.SessionUpdate(st)
	if v, ok := st.Data[key]; ok {
		return v
	} else {
		return nil
	}
	return nil
}

//删除
func (st *SessionStore) Delete(key string) error {
	delete(st.Data, key)
	redisProvider.SessionUpdate(st)
	return nil
}

//获取sessionID
func (st *SessionStore) SessionID() string {
	return st.Sid
}

func (rp *RedisProvider) SessionInit(sid string) (session.Session, error) {
	if sid == "" {
		return nil, fmt.Errorf("sid 为空")
	}
	rp.lock.Lock()
	defer rp.lock.Unlock()
	v := make(map[string]interface{}, 0)
	timestamp := time.Now()
	st := &SessionStore{
	    Sid: sid,
	    LastAccessedTime: timestamp,
	    Data: v,
	}
	rp.stNumber[sid] = timestamp.Unix()
	err := rp.SessionUpdate(st)
	return st, err
}

func (rp *RedisProvider) parseSessionStore(sid string) (*SessionStore, error) {
    st := &SessionStore{}
    var err error
    var element string
    if element, err = rp.Client.Get(sid).Result(); err != nil {
		if err == redis.Nil {
			return nil, nil
		} else {
			return nil, err
		}
    }
	err = json.Unmarshal([]byte(element), st)
	if err != nil {
		return nil, fmt.Errorf("[session] 解析失败, %v", err)
	}
	return st, nil

}

func (rp *RedisProvider) SessionRead(sid string) (session.Session, error) {
	st, err := rp.parseSessionStore(sid)
	if st != nil {
		rp.SessionUpdate(st)
		return st, err
	} else {
		//特别注意, 因为`rp.parseSessionStore(sid)`返回的是*SessionStore类型, 即使*SessionStore = nil, 但 Session != nil, 因为Session的动态类型为*SessionStore
		return nil, err
	}
}

func (rp *RedisProvider) SessionDestroy(sid string) error {
    _, err := rp.Client.Del(sid).Result()
    return err
}

func (rp *RedisProvider) SessionGC(maxLifeTime int64) int {
	rp.lock.Lock()
	defer rp.lock.Unlock()
	//对话会话进行排序
    result := common.SortMapSI64(rp.stNumber)
    //遍历排序后的结果, 并清除对应的会话
    lastSession := len(result)
	for _, sid := range result {
	    if _, ok := rp.stNumber[sid]; ok {
            session, err := rp.parseSessionStore(sid)
            if err != nil {
                //todo log
                continue
            }
            //经过排序清楚过期会话
            if session.LastAccessedTime.Unix() < time.Now().Unix() {
				session.Delete(sid)
				lastSession--
            } else {
                break;
            }
        }
    }
	return lastSession
}
func (rp *RedisProvider) SessionUpdate(session *SessionStore) error {
	session.LastAccessedTime = time.Now()
    jsonSt, err := json.Marshal(session)
    if err != nil {
        return fmt.Errorf("[session] SessionStore转换为json失败, %v", err)
    }
    _, err = rp.Client.Set(session.Sid, string(jsonSt), 2 * time.Hour).Result();
	if err != nil {
		return fmt.Errorf("[session] SessionStore缓存到redis失败, %v", err)
	}
	return nil
}
