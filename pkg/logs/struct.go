package logs

import (
	"log"
	"os"
)

// 日志级别
type logLevel int32

// 日志结构
type logger struct {
	LogLevel int32 // 日志级别

	TraceMessage   *log.Logger // 追踪信息
	InfoMessage    *log.Logger // 普通信息
	WarningMessage *log.Logger // 警告信息
	ErrorMessage   *log.Logger // 错误信息
	FatalMessage   *log.Logger // 致命信息

	OutHandler *handler // 输出对象
}

// 输出对象
type handler struct {
	fs   *os.File
	name string
}
