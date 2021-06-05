package util

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// FailFile 将整理失败的文件存储到fail目录中
//
// file 字符串参数，传入失败文件路径，
// fail 字符串参数，传入fail目录路径。
func FailFile(file, fail string) {
	// 获取运行路径
	base := GetRunPath()
	// 组合路径
	base = base + "/" + fail

	// 移动文件到失败目录
	err := MoveFile(file, base+"/"+path.Base(file))
	// 检查
	if err != nil {
		return
	}
}

// MoveFile 移动文件到指定路径，并返回错误信息
//
// oldPath 字符串参数，传入文件原始路径，
// newPath 字符串参数，传入文件移动路径。
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

// GetFileSize 获取指定文件大小，失败则返回0
//
// file 字符串参数，传入文件路径
func GetFileSize(file string) int64 {
	// 获取文件信息
	info, err := os.Stat(file)
	if err != nil {
		return 0
	}

	// 获取大小
	return info.Size()
}

// WriteFile 将字节集数据写入到指定文件中，并返回错误信息
//
// file 字符串参数，传入写入文件路径，
// data 字节集参数，传入写入的数据。
func WriteFile(file string, data []byte) error {
	// 写文件
	return ioutil.WriteFile(file, data, 0644)
}

// ReadFile 读取文件
//
// file 字符串参数，传入文件路径
func ReadFile(file string) ([]byte, error) {
	b, err := ioutil.ReadFile(file)

	if err != nil {
		return nil, err
	}

	return b, nil
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
