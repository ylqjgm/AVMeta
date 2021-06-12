package logs

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync/atomic"
)

// Log 创建日志对象
func Log(file string) {
	// 级别信息
	var traceHandler, infoHandler, warnHandler io.Writer = os.Stdout, os.Stdout, os.Stdout
	var errorHandler, fatalHandler io.Writer = os.Stderr, os.Stderr

	// 创建日志文件
	handler, err := initHandler(file)
	if err != nil {
		panic(err)
	}

	// 是否传入文件名
	if file != "" {
		// 写入追踪
		if traceHandler == os.Stdout {
			traceHandler = io.MultiWriter(handler, traceHandler)
		}

		// 写入普通信息
		if infoHandler == os.Stdout {
			infoHandler = io.MultiWriter(handler, infoHandler)
		}

		// 写入警告信息
		if warnHandler == os.Stdout {
			warnHandler = io.MultiWriter(handler, warnHandler)
		}

		// 写入错误信息
		if errorHandler == os.Stderr {
			errorHandler = io.MultiWriter(handler, errorHandler)
		}

		// 写入致命信息
		if fatalHandler == os.Stderr {
			fatalHandler = io.MultiWriter(handler, fatalHandler)
		}
	}

	// 初始化
	Logger = &logger{
		TraceMessage:   log.New(traceHandler, "[Trace]: ", Flags),
		InfoMessage:    log.New(infoHandler, "[Info]: ", Flags),
		WarningMessage: log.New(warnHandler, "[Warning]: ", Flags),
		ErrorMessage:   log.New(errorHandler, "[Error]: ", Flags),
		FatalMessage:   log.New(fatalHandler, "[Fatal]: ", Flags),
		OutHandler:     handler,
	}

	// 内存存储
	atomic.StoreInt32(&Logger.LogLevel, int32(LevelTrace|LevelInfo|LevelWarn|LevelError))
}

// Close 关闭日志
func Close() {
	if Logger.OutHandler != nil {
		_ = Logger.OutHandler.Close()
	}
}

// Sync 日志同步
func sync() {
	if Logger.OutHandler != nil {
		_ = Logger.OutHandler.fs.Sync()
	}
}

// Trace 追踪日志
func Trace(format string, a ...interface{}) {
	_ = Logger.TraceMessage.Output(2, fmt.Sprintf(format, a...))
}

// Info 普通信息
func Info(format string, a ...interface{}) {
	_ = Logger.InfoMessage.Output(2, fmt.Sprintf(format, a...))
}

// Warning 警告信息
func Warning(format string, a ...interface{}) {
	_ = Logger.WarningMessage.Output(2, fmt.Sprintf(format, a...))
}

// Error 错误信息
func Error(format string, a ...interface{}) {
	_ = Logger.ErrorMessage.Output(2, fmt.Sprintf(format, a...))
}

// Fatal 按照格式写入致命信息
func Fatal(format string, a ...interface{}) {
	// 写入日志
	_ = Logger.FatalMessage.Output(2, fmt.Sprintf(format, a...))
	// 日志同步
	sync()
	// 关闭日志
	Close()
	// 退出执行
	os.Exit(255)
}

// FatalError 出错则写入致命错误
func FatalError(err error) {
	// 检查错误
	if err != nil {
		Fatal("%s\n", err)
	}
}
