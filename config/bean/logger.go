package bean

import "time"

type LogLevel string

var (
	Debug LogLevel = "debug"
	Info  LogLevel = "info"
	Warn  LogLevel = "warn"
	Error LogLevel = "error"
	Fatal LogLevel = "fatal"
	Panic LogLevel = "panic"
)

type LogFile struct {
	FileNameFormat string     `json:"file_name_format"`         //日志文件名标准格式
	Level          []LogLevel `json:"levels"`                   //不同级别日志路由
	LinkName       string     `json:"link_file_name omitempty"` //软连接, ""空表示不使用
}

type Logger struct {
	RotationTime  time.Duration `json:"rotation_time"` //分割文件周期, 单位s
	RotationCount uint          `json:"rotation_cnt"`  //最大日志文件个数
	StandLevel    LogLevel      `json:"stand_level"`   //监听日志日志级别下限
	Files         []LogFile     `json:"files"`         //日志文件属性
}
