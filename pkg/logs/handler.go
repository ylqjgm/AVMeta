package logs

import (
	"github.com/ylqjgm/AVMeta/pkg/util"
	"os"
	"time"
)

// 初始化输出对象
func initHandler(file string) (*handler, error) {
	// 是否传入文件
	if file == "" {
		return &handler{fs: nil, name: ""}, nil
	}
	// 日志路径
	filePath := util.GetRunPath() + "/log/" + time.Now().Format("20060102") + "/"
	// 创建目录
	err := os.MkdirAll(filePath, 0777)
	// 检查错误
	if err != nil {
		return nil, err
	}
	// 创建对象
	h := new(handler)
	h.name = filePath + file + ".log"

	// 创建并打开文件
	h.fs, err = os.OpenFile(h.name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return h, nil
}

// Write 写入日志
func (h *handler) Write(b []byte) (int, error) {
	return h.fs.Write(b)
}

// Close 关闭日志
func (h *handler) Close() error {
	if h.fs != nil {
		return h.fs.Close()
	}

	return nil
}
