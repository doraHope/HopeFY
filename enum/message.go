package enum

var errMsg = map[int]string{
	SUCCESS:        "成功",
	INVALID_PARAMS: "参数错误",
	SERVICE_ERROR:  "服务异常",
}

func ErrMsg(code int) string {
	if msg, ok := errMsg[code]; ok {
		return msg
	}
	return errMsg[SERVICE_ERROR]
}
