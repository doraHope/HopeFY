package middlerware

import (
	"encoding/json"
	"sync"
	"time"

	//工具
	"../util"
	//配置
	"../config"
)

const (
	GC_INTERVAL = 30 //30min
)

func init() {
	mg = Manager{
		sm:             make(map[string]int, 32),
		sl:             make([]SessionStore, 0, 32),
		lock:           sync.RWMutex{},
		prefix:         config.AppConfig.Session.Prefix,
		expired:        config.AppConfig.Session.Expired,
		maxOpenSession: config.AppConfig.Session.MaxOpenSession,
	}
	Cont.Register("session", &mg)
	go func() {
		for {
			select {
			case <-time.After(GC_INTERVAL * time.Minute):
				//计时
				ts := time.Now()
				mg.gc()
				tc := time.Since(ts).Milliseconds()
				//todo log tc
				_ = tc
			}
		}
	}()
}

func GetSession(s interface{}) *Manager {
	m, ok := s.(*Manager)
	if !ok {
		return nil
	}
	return m
}

type SessionError struct {
	err string
}

func (s *SessionError) Error() string {
	return s.err
}

var SessionNil = &SessionError{
	"session is nil",
}

var SessionNotSupportType = &SessionError{
	"value数据类型, 仅支持json所支持的数据类型",
}

var SessionNotFound = &SessionError{
	"未找到key关联的value",
}

var SessionNotInit = &SessionError{
	"session还未初始化",
}

var ManagerNotInit = &SessionError{
	"manager 还未初初始化",
}

type SessionStore struct {
	id       string //sessionID
	data     map[string]interface{}
	expireAt time.Time //过期时间
	lastUpAt time.Time //上一次更新时间
}

type Manager struct {
	sm             map[string]int //索引表
	sl             []SessionStore //数组结构
	lock           sync.RWMutex   //互斥锁
	prefix         string         //前缀
	expired        int            //有效时间, 单位m
	openSession    int            //当前管理的会话数量
	maxOpenSession int            //最大会话数量
}

var mg Manager

func (s *SessionStore) Set(key string, value interface{}) error {
	switch v := value.(type) {
	case int8:
		s.data[key] = int64(v)
	case int16:
		s.data[key] = int64(v)
	case int32:
		s.data[key] = int64(v)
	case int:
		s.data[key] = int64(v)
	case int64:
		s.data[key] = v
	case uint8:
		s.data[key] = uint64(v)
	case uint16:
		s.data[key] = uint64(v)
	case uint32:
		s.data[key] = uint64(v)
	case uint:
		s.data[key] = uint64(v)
	case uint64:
		s.data[key] = v
	case float32:
		s.data[key] = float64(v)
	case float64:
		s.data[key] = v
	case string:
		s.data[key] = v
	case bool:
		if v {
			s.data[key] = true
		} else {
			s.data[key] = false
		}
	default:
		return SessionNotSupportType
	}
	s.update()
	return nil
}

func (s *SessionStore) Get(key string) interface{} {
	if s == nil {
		return nil
	}
	if v, ok := s.data[key]; ok {
		return v
	}
	return nil
}

func (s *SessionStore) Delete(key string) error {
	if s == nil {
		return SessionNotInit
	}
	if _, ok := s.data[key]; ok {
		delete(s.data, key)
	} else {
		return SessionNotFound
	}
	//do what
	return nil
}

func (s *SessionStore) SessionID() (string, error) {
	if s == nil {
		return "", SessionNotInit
	}
	return s.id, nil
}

func (s *SessionStore) clean() error {
	//置为过期
	s.expireAt = time.Unix(0, 0)
	s.lastUpAt = time.Unix(0, 0)
	//清空缓存
	_, err := rc.Del(s.id).Result()
	if err != nil {
		//todo log
	}
	return nil
}

//数据同步到缓存中
func (s *SessionStore) update() {
	if len(s.data) > 0 {
		//sync cache
		jd, err := json.Marshal(s.data)
		if err != nil {
			//todo log
		}
		_, err = rc.Set(s.id, jd, time.Duration(mg.expired)*time.Minute).Result()
		if err != nil {
			//todo log
		}
		s.lastUpAt = time.Now()
		s.expireAt = time.Now().Add(time.Duration(mg.expired) * time.Minute)
	}
}

//会话初始化
func (m *Manager) SessionInit() (*SessionStore, error) {
	if m == nil {
		return nil, ManagerNotInit
	}
	sid := util.CreateUUID()
	if sid == "" {
		return nil, SessionNotInit
	}
	sid = m.prefix + "_" + sid
	s, _ := m.getSS(sid)
	if s == nil {
		s = m.setSS(sid)
	}
	s.update()
	return s, nil
}

func (m *Manager) SessionRead(sid string) (*SessionStore, error) {
	if m == nil {
		return nil, ManagerNotInit
	}
	s, up := m.getSS(sid)
	if s == nil {
		return nil, SessionNotFound
	}
	if up {
		s.update()
		m.heapDown(m.sm[s.id])
	}
	return s, nil
}

//销毁会话
func (m *Manager) SessionDestroy(sid string) error {
	//仅清空索引表对应的项
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.sm[sid]; ok {
		delete(m.sm, sid)
		return nil
	}
	return SessionNotFound
}

//清除过期会话-暂时不考虑会话上限
func (m *Manager) gc() {
	m.lock.Lock()
	defer m.lock.Unlock()
	if len(m.sl) > 0 {
		for {
			s := m.sl[0]
			timestamp := time.Now().Unix()
			if s.expireAt.Unix() < timestamp {
				//清空session缓存中的内容
				s.clean()
				li := len(m.sl) - 1
				//从索引中删除该session
				delete(m.sm, s.id)
				//堆删除
				m.sl[0], m.sl[li] = m.sl[li], m.sl[0]
				m.sm[m.sl[0].id] = 0
				//如果此时堆的大小为0, 直接赋值nil
				if li == 0 {
					m.sl = nil
				} else {
					//堆裁剪
					m.sl = m.sl[:li-1]
					m.openSession--
					//堆重构
					m.heapDown(0)
				}
			}
			break
		}
	}
}

//堆重构-下沉
func (m *Manager) heapDown(i int) {
	hl := len(m.sl) - 1
	for {
		nl := i*2 + 1
		nr := i*2 + 2
		if nl <= hl && nr <= hl {
			//与左节点比较, 如果左右节点相等的情况, 左子节点优先进行
			if m.sl[nl].lastUpAt.Unix() >= m.sl[nr].lastUpAt.Unix() && m.sl[i].lastUpAt.Unix() < m.sl[nl].lastUpAt.Unix() {
				m.sm[m.sl[nl].id], m.sm[m.sl[i].id] = m.sm[m.sl[i].id], m.sm[m.sl[nl].id]
				m.sl[nl], m.sl[i] = m.sl[i], m.sl[nl]
				i = nl
				continue
				//与右节点比较
			} else if m.sl[i].lastUpAt.Unix() < m.sl[nr].lastUpAt.Unix() {
				m.sm[m.sl[nr].id], m.sm[m.sl[i].id] = m.sm[m.sl[i].id], m.sm[m.sl[nr].id]
				m.sl[nr], m.sl[i] = m.sl[i], m.sl[nr]
				i = nr
				continue
			}
		} else if nl <= hl && m.sl[i].lastUpAt.Unix() < m.sl[nl].lastUpAt.Unix() {
			//同末尾节点比较
			m.sm[m.sl[nl].id], m.sm[m.sl[i].id] = m.sm[m.sl[i].id], m.sm[m.sl[nl].id]
			m.sl[nl], m.sl[i] = m.sl[i], m.sl[nl]
		}
		break
	}
}

//堆重构-上浮
func (m *Manager) heapUp() int {
	//上沉
	i := len(m.sl) - 1
	for {
		pi := (i - 1) / 2
		if m.sl[i].lastUpAt.Unix() < m.sl[pi].lastUpAt.Unix() {
			m.sl[i], m.sl[pi] = m.sl[pi], m.sl[i]
			i = pi
			continue
		}
		break
	}
	return i
}

//获取会话
func (m *Manager) getSS(sid string) (*SessionStore, bool) {
	m.lock.RLock()
	if i, ok := m.sm[sid]; ok {
		m.lock.RUnlock()
		return &m.sl[i], true
	} else {
		m.lock.RUnlock()
		//尝试从缓存中间件中读取
		sc := make(map[string]interface{})
		if sd, err := rc.Get(sid).Result(); err == nil {
			err := json.Unmarshal([]byte(sd), &sc)
			if err == nil {
				m.lock.Lock()
				defer m.lock.Unlock()

				m.sl = append(m.sl, SessionStore{
					id:       sid,
					data:     sc,
					expireAt: time.Now().Add(time.Duration(m.expired) * time.Minute),
					lastUpAt: time.Now(),
				})
				m.sm[sid] = m.openSession
				m.openSession++
				return &m.sl[m.sm[sid]], false
			} else {
				//todo log
			}
		}
	}
	return nil, false
}

//设置会话
func (m *Manager) setSS(sid string) *SessionStore {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.sl = append(m.sl, SessionStore{
		id:       sid,
		data:     make(map[string]interface{}, 8),
		expireAt: time.Now().Add(time.Duration(m.expired) * time.Minute),
		lastUpAt: time.Now(),
	})
	m.sm[sid] = m.openSession
	m.openSession++
	return &m.sl[m.sm[sid]]
}
