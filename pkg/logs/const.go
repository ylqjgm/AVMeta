package logs

import "log"

const (
	LevelTrace logLevel                                 = 1 << iota // 追踪日志
	LevelInfo                                                       // 普通信息
	LevelWarn                                                       // 警告信息
	LevelError                                                      // 错误信息
	Flags      = log.Ldate | log.Ltime | log.Lshortfile             // 日志参数
)

var Logger *logger
