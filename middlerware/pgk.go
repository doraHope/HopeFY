//包装方法
package middlerware

//启动session
func SessionStart(sid *string) *SessionStore {
	var mg *Manager
	var ss *SessionStore

	s, err := Cont.Get("session")
	if err != nil {
		//todo log
		return nil
	}
	mg = GetSession(s)
	if sid == nil || *sid == "" {
		ss, err = mg.SessionInit()
	} else {
		ss, err = mg.SessionRead(*sid)
	}
	if err != nil {
		//todo log
		return nil
	}
	return ss
}

func Log(msg, lv string) {
    switch lv {
    case "debug":
        lg.Debug(msg)
    case "info":
        lg.Info(msg)
    case "warn":
        lg.Warn(msg)
    case "error":
        lg.Error(msg)
    case "fatal":
        lg.Fatal(msg)
    case "panic":
        lg.Panic(msg)
    default:
        lg.Warnf("[util]-`错误日志类型`, lv:`%v`", lv)
    }
}
