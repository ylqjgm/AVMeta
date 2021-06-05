package util

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// NfoFile nfo文件结构体
type NfoFile struct {
	Path   string
	Fanart string
	Poster string
	Dir    string
}

// WalkDir 遍历指定目录下的视频文件，
// 返回文件路径列表及错误信息。
//
// dirPath 字符串参数，传入要遍历的目录路径，
// success 字符串参数，传入要过滤的success目录名称，
// fail 字符串参数，传入要过滤的fail目录名称。
func WalkDir(dirPath, success, fail string) ([]string, error) {
	// 定义文件列表
	var files []string

	// 遍历目录
	err := filepath.Walk(dirPath, func(filePath string, f os.FileInfo, err error) error {
		// 错误
		if f == nil {
			return err
		}
		// 忽略目录
		if f.IsDir() {
			return nil
		}

		// 检测是否为过滤目录
		if strings.Contains(
			strings.ToUpper(filePath),
			strings.ToUpper(success)) ||
			strings.Contains(
				strings.ToUpper(filePath),
				strings.ToUpper(fail)) {
			return nil
		}

		// 隐藏文件正则
		rHidden := regexp.MustCompile(`^\.(.)*`)
		// 检测是否为隐藏文件
		if rHidden.MatchString(f.Name()) {
			return nil
		}

		// 获取后缀并转换为小写
		ext := strings.ToLower(path.Ext(filePath))

		// 验证是否存在于后缀扩展名中
		if _, ok := videoExts[ext]; ok {
			// 存在则加入扩展
			files = append(files, filePath)
		}

		return nil
	})

	return files, err
}

// GetRunPath 获取程序当前执行路径
func GetRunPath() string {
	// 获取当前执行路径
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	// 检查错误
	if err != nil {
		return "."
	}

	return dir
}

// WalkNfo 遍历当前目录下所有 .nfo 文件
func WalkNfo(dirPath string, files []NfoFile) ([]NfoFile, error) {
	// 读取目录
	r, err := ioutil.ReadDir(dirPath)
	// 检查错误
	if err != nil {
		return nil, err
	}

	// 循环列表
	for _, f := range r {
		if f.IsDir() {
			fullDir := dirPath + "/" + f.Name()
			files, err = WalkNfo(fullDir, files)
			if err != nil {
				return files, err
			}
		} else {
			if strings.ToLower(filepath.Ext(f.Name())) == ".nfo" {
				// nfo变量
				var nfo NfoFile
				// nfo路径
				nfo.Path = dirPath + "/" + f.Name()
				// nfo目录
				nfo.Dir = dirPath
				// 是否存在fanart.jpg
				if Exists(dirPath + "/fanart.jpg") {
					nfo.Fanart = dirPath + "/fanart.jpg"
				}
				// 是否存在poster.jpg
				if Exists(dirPath + "/poster.jpg") {
					nfo.Poster = dirPath + "/poster.jpg"
				}

				// 加入文件列表
				files = append(files, nfo)
			}
		}
	}

	return files, nil
}
