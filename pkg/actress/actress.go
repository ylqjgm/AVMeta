package actress

import (
	"fmt"
	"github.com/ylqjgm/AVMeta/pkg/logs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/ylqjgm/AVMeta/pkg/util"

	"github.com/schollz/progressbar/v2"
)

const (
	// JAVBUS 是 javbus 网站名称常量
	JAVBUS = "JAVBUS"
	// JAVDB 是 javdb 网站名称常量
	JAVDB = "JAVDB"
)

// Actress 头像管理结构体
type Actress struct {
	// 程序配置
	cfg *util.ConfigStruct
	// Emby 媒体库API对象
	emby *Emby
}

// NewActress 返回一个Actress对象。
// 引用util.GetConfig方法读取配置文件，
// 配置文件读取失败则返回空对象。
func NewActress() *Actress {
	// 初始化日志
	logs.Log("actress")

	// 获取配置信息
	cfg, err := util.GetConfig()
	// 检查
	logs.FatalError(err)

	return &Actress{
		cfg:  cfg,
		emby: NewEmby(cfg.Media.URL, cfg.Media.API),
	}
}

// Fetch 远程女优头像下载。
// 通过传入参数获取远程网站女优头像图片并下载到本地。
// 所有图片均下载到程序执行目录下的 actress 文件夹中。
//
// site 字符串参数，指定要下载的网站名称，参见常量定义，
// page 整数参数，指定要下载的开始页面，
// censored 逻辑参数，指定下载的是有码女优还是无码女优。
func (a *Actress) Fetch(site string, page int, censored bool) error {
	// 定义数据存储map
	var acts map[string]string
	// 定义下一页变量
	var next bool
	// 定义错误变量
	var err error

	// 根据不同的站点选择不同的处理方式
	switch site {
	case JAVBUS: // javBUS
		// 采集
		acts, next, err = JavBUS(a.cfg.Site.JavBus, a.cfg.Base.Proxy, page, censored)
	case JAVDB: // javDB
		// 采集
		acts, next, err = JavDB(a.cfg.Site.JavDB, a.cfg.Base.Proxy, page, censored)
	default:
		return fmt.Errorf("site case error")
	}
	// 检查
	if err != nil {
		return err
	}

	// 总量
	count := len(acts)

	// 初始化进程
	wg := util.NewWaitGroup(5)
	// 定义进度条
	bar := progressbar.NewOptions(count,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() { fmt.Println("") }),
		progressbar.OptionSetDescription(fmt.Sprintf("第 [blue][%d][reset] 页, [green][%d][reset] 位女优...", page, count)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[cyan]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[cyan][[reset]",
			BarEnd:        "[cyan]][reset]",
		}),
	)

	// 循环下载
	for name, cover := range acts {
		// 计数加
		wg.AddDelta()
		// 调用
		go a.downProcess(name, cover, wg, bar)
	}

	// 等待结束
	wg.Wait()

	// 不管有没有全部下载，均设置为已完成
	_ = bar.Finish()

	if next {
		// 采集下一页
		return a.Fetch(site, page+1, censored)
	}

	return nil
}

// Put 本地图片入库
// 扫描程序执行目录下的 actress 文件夹，
// 将其中的所有女优头像依次入库到 Emby 中。
func (a *Actress) Put() error {
	// 获取文件列表
	files, err := a.walkDir()
	// 检查
	if err != nil {
		logs.Error("获取头像列表失败, 错误信息: %s\n", err)
		return err
	}

	// 获取总量
	count := len(files)
	// 进度条
	bar := progressbar.NewOptions(count,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionOnCompletion(func() { fmt.Println("") }),
		progressbar.OptionSetDescription(fmt.Sprintf("[blue][%d/%d][reset]...", 0, count)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[cyan]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[cyan][[reset]",
			BarEnd:        "[cyan]][reset]",
		}),
	)

	// 循环文件
	for k, f := range files {
		// 限制速度
		time.Sleep(300 * time.Millisecond)
		// 获取后缀
		ext := path.Ext(f)
		// 提取女优名称
		name := strings.TrimRight(path.Base(f), ext)

		bar.Describe(fmt.Sprintf("[cyan][%d/%d][reset]", k+1, count))
		_ = bar.Set(k + 1)

		// 调用头像上传
		err = a.emby.Actor(name, f)
		// 检查
		if err != nil {
			continue
		}

		// 标注成功
		a.success(f)
	}

	_ = bar.Finish()

	return nil
}

// 下载多进程处理
func (a *Actress) downProcess(name, cover string, wg *util.WaitGroup, bar *progressbar.ProgressBar) {
	// 检测是否已存在或入库过
	if a.exists(name) {
		wg.Done()
		_ = bar.Add(1)
		return
	}
	// 获取扩展
	ext := path.Ext(cover)
	// 下载图片
	_ = util.SavePhoto(cover,
		fmt.Sprintf("./actress/%s.jpg", name),
		a.cfg.Base.Proxy,
		!strings.EqualFold(strings.ToLower(ext), ".jpg"))

	wg.Done()
	_ = bar.Add(1)
}

// 女优文件是否存在
func (a *Actress) exists(name string) bool {
	// 获取文件信息
	_, err := os.Stat("./actress/" + name + ".jpg")
	// 检查错误
	if err == nil {
		return true
	}
	// 是否不存在
	if os.IsNotExist(err) {
		return false
	}

	// 获取文件信息
	_, err = os.Stat("./actress/success/" + name + ".jpg")
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

// 遍历头像目录
func (a *Actress) walkDir() (files []string, err error) {
	// 遍历目录
	err = filepath.Walk("./actress", func(filePath string, f os.FileInfo, err error) error {
		// 错误
		if f == nil {
			return err
		}
		// 忽略目录
		if f.IsDir() {
			return nil
		}

		// 是否为成功的
		if strings.EqualFold(path.Base(path.Dir(filePath)), "success") {
			return nil
		}

		// 隐藏文件正则
		rHidden := regexp.MustCompile(`^\.(.)*`)
		// 检测是否为隐藏文件
		if rHidden.MatchString(f.Name()) {
			return nil
		}
		// 加入列表
		files = append(files, filePath)

		return nil
	})

	return
}

// 标记已入库
func (a *Actress) success(file string) {
	// 获取文件名
	fname := path.Base(file)
	// 设定success路径
	successDir := "./actress/success"
	// 创建success目录
	_ = os.MkdirAll(successDir, os.ModePerm)
	// 移动文件
	_ = os.Rename("./actress/"+fname, successDir+"/"+fname)
}
