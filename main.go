package main

import (
	//程序初始化
	_ "./init"

	//测试
	mw "./middlerware"
)

func main() {
	//session 测试
	//var sid *string = new(string)
	//*sid = `swy_VN7lqXWXMU2I8gQgbU0uXzfg_ryCcTn2uKClR_lVRUU=`
	//cs := middlerware.SessionStart(sid)
	////cs.Set("ceshi", "hope for you!")
	//data := cs.Get("ceshi")
	//if data != nil {
	//	fmt.Println("data", data)
	//}

	//logger测试
	//li, err := middlerware.Cont.Get("log")
	//if err == nil {
	//	lg := middlerware.GetLog(li)
	//	if lg != nil {
	//		lg.Info(`hello world`)
	//	} else {
	//		fmt.Println(`error null`)
	//	}
	//}

	mw.Log(`hello world`, `info`)

}
