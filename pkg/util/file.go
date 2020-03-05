package util

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// FailFile 恢复文件
func FailFile(file string, cfg *ConfigStruct) {
	// 获取运行路径
	base := GetRunPath()
	// 组合路径
	base = base + "/" + cfg.Path.Fail

	// 移动文件到失败目录
	err := MoveFile(file, base+"/"+path.Base(file))
	// 检查
	if err != nil {
		return
	}
}

// MoveFile 移动文件
func MoveFile(oldPath, newPath string) error {
	// 创建目录
	err := os.MkdirAll(filepath.Dir(newPath), os.ModePerm)
	// 检查错误
	if err != nil {
		return err
	}
	// 移动文件
	return os.Rename(oldPath, newPath)
}

// GetFileSize 获取文件大小
func GetFileSize(file string) int64 {
	// 获取文件信息
	info, err := os.Stat(file)
	if err != nil {
		return 0
	}

	// 获取大小
	return info.Size()
}

// WriteFile 写入文件
func WriteFile(file string, data []byte) error {
	// 写文件
	return ioutil.WriteFile(file, data, 0644)
}

// Exists 检查文件是否存在
func Exists(file string) bool {
	// 获取文件信息
	_, err := os.Stat(file)
	// 检查错误
	if err == nil {
		return true
	}
	// 是否不存在
	if os.IsNotExist(err) {
		return false
	}

	return false
}
