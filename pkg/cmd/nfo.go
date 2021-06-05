package cmd

import (
	"encoding/xml"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/ylqjgm/AVMeta/pkg/media"
	"github.com/ylqjgm/AVMeta/pkg/util"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	// 定义扩展列表
	videoExts = map[string]string{
		".avi":  ".avi",
		".flv":  ".flv",
		".mkv":  ".mkv",
		".mov":  ".mov",
		".mp4":  ".mp4",
		".rmvb": ".rmvb",
		".ts":   ".ts",
		".wmv":  ".wmv",
	}
)

// NfoFile nfo文件列表结构
type NfoFile struct {
	Path  string
	Video string
	Dir   string
}

// nfo命令
func (e *Executor) initNfo() {
	nfoCmd := &cobra.Command{
		Use: "nfo",
		Long: `
自动将运行目录下所有nfo文件转换为VSMeta文件`,
		Example: `  AVMeta nfo`,
		Run:     e.nfoRunFunc,
	}

	e.rootCmd.AddCommand(nfoCmd)
}

// 转换执行命令
func (e *Executor) nfoRunFunc(cmd *cobra.Command, args []string) {
	// 获取当前执行路径
	curDir := util.GetRunPath()

	// 文件列表
	var nfos []NfoFile

	// 列当前目录
	nfos, err := e.walk(curDir, nfos)
	// 检测错误
	if err != nil {
		log.Fatalln(err)
	}

	// 获取总量
	count := len(nfos)
	// 输出总量
	fmt.Printf("\n共探索到 %d 个 nfo 文件, 开始转换...\n\n", count)

	// 初始化进程
	wg := util.NewWaitGroup(2)

	// 循环nfo文件列表
	for _, nfo := range nfos {
		// 计数加
		wg.AddDelta()
		// 转换进程
		go e.nfoProcess(nfo, wg)
	}

	// 等待结束
	wg.Wait()
}

// 转换进程
func (e *Executor) nfoProcess(nfo NfoFile, wg *util.WaitGroup) {
	// 读取文件
	b, err := util.ReadFile(nfo.Path)
	// 检查
	if err != nil {
		// 输出错误
		fmt.Printf("文件: [%s] 打开失败, 错误原因: %s\n", path.Base(nfo.Path), err)

		// 进程
		wg.Done()

		return
	}

	// 媒体对象
	var m media.Media

	// 转换
	err = xml.Unmarshal(b, &m)
	// 检查错误
	if err != nil {
		// 输出错误
		fmt.Printf("文件: [%s] 打开失败, 错误原因: %s\n", path.Base(nfo.Path), err)

		// 进程
		wg.Done()

		return
	}

	// scheme
	var scheme = map[string]string{
		"http":  "http",
		"https": "https",
	}

	// 实例化vsmeta
	vs := media.NewVSMeta()
	// fanart
	if util.Exists(nfo.Dir + "/fanart.jpg") {
		m.FanArt = nfo.Dir + "/fanart.jpg"
	} else if m.FanArt != "" {
		uri, err := url.Parse(m.FanArt)
		if err != nil || uri == nil {
			m.FanArt = ""
		}
		if _, ok := scheme[uri.Scheme]; !ok {
			m.FanArt = ""
		}
	}
	// poster
	if util.Exists(nfo.Dir + "/poster.jpg") {
		m.Poster = nfo.Dir + "/poster.jpg"
	} else if m.Poster != "" {
		uri, err := url.Parse(m.Poster)
		if err != nil || uri == nil {
			m.Poster = ""
		}
		if _, ok := scheme[uri.Scheme]; !ok {
			m.Poster = ""
		}
	}

	// 解析为 vsmeta
	bs := vs.Convert(&m)

	// 获取视频后缀
	ext := path.Ext(nfo.Video)

	// 写入vsmeta
	err = util.WriteFile(fmt.Sprintf("%s/%s%s.vsmeta", nfo.Dir, m.Number, ext), bs)
	// 检查
	if err != nil {
		// 输出错误
		fmt.Printf("文件: [%s] 转换失败, 错误原因: %s\n", path.Base(nfo.Path), err)

		// 进程
		wg.Done()

		return
	}

	// 输出正确
	fmt.Printf("文件: [%s/%s] 转换成功, 路径: %s\n", path.Base(nfo.Path), m.Number, nfo.Dir)

	// 进程
	wg.Done()
}

// 列目录
func (e *Executor) walk(dirPath string, nfoFiles []NfoFile) ([]NfoFile, error) {
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
			nfoFiles, err = e.walk(fullDir, nfoFiles)
		} else {
			// 获取后缀并转换为小写
			ext := strings.ToLower(path.Ext(f.Name()))

			// 验证是否存在于后缀扩展名中
			if _, ok := videoExts[ext]; ok {
				// 遍历目录
				err := filepath.Walk(dirPath, func(filePath string, fi os.FileInfo, err error) error {
					// 错误
					if fi == nil {
						return err
					}

					if !fi.IsDir() {
						// 是否nfo文件
						if strings.ToLower(filepath.Ext(fi.Name())) == ".nfo" {
							// 初始化nfo
							nfo := NfoFile{
								Path:  dirPath + "/" + fi.Name(),
								Video: dirPath + "/" + f.Name(),
								Dir:   dirPath,
							}

							// 加入列表
							nfoFiles = append(nfoFiles, nfo)
						}
					}

					return nil
				})

				if err != nil {
					return nfoFiles, err
				}
			}
		}
	}

	return nfoFiles, nil
}
