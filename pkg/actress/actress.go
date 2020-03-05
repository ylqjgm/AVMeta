package actress

import (
	"fmt"
	"path"
	"strings"

	"github.com/ylqjgm/AVMeta/pkg/util"

	"github.com/PuerkitoBio/goquery"
)

// Actress 头像对象
type Actress struct {
	JavBUS string
	JavDB  string
	Media  string
	cfg    *util.ConfigStruct
}

// NewActress 创建对象
func NewActress(javbus, javdb, media string) *Actress {
	// 获取配置信息
	cfg, err := util.GetConfig()
	// 检查
	if err != nil {
		return nil
	}

	return &Actress{
		JavBUS: javbus,
		JavDB:  javdb,
		Media:  media,
		cfg:    cfg,
	}
}

// Search 搜索女优头像
func (a *Actress) Search(name string) string {
	// 定义头像变量
	var actress string

	// 先从javBus搜索
	actress = a.javBusSearch(name)
	// 检查
	if actress == "" {
		// 从javDB搜索
		actress = a.javDBSearch(name)
	}

	return actress
}

// Save 保存头像
func (a *Actress) Save(name string) error {
	// 保存路径
	savePath := fmt.Sprintf("%s/actress/%s.jpg", util.GetRunPath(), name)
	// 检查头像是否已经存在了
	if util.Exists(savePath) {
		return nil
	}

	// 先搜索到头像
	actress := a.Search(name)
	// 检查下
	if actress == "" {
		return fmt.Errorf("404 Not Found")
	}

	// 获取头像后缀
	ext := path.Ext(actress)
	// 检查是否需要转换
	if !strings.EqualFold(strings.ToLower(ext), ".jpg") {
		return util.SavePhoto(actress, savePath, a.cfg.Base.Proxy, true)
	}

	return util.SavePhoto(actress, savePath, a.cfg.Base.Proxy, false)
}

// Stock 头像入库
func (a *Actress) Stock(name string) error {
	// 先保存头像
	if err := a.Save(name); err != nil {
		return err
	}

	// 头像路径
	acctressPath := fmt.Sprintf("%s/actress/%s.jpg", util.GetRunPath(), name)

	// 检查是否配置了api
	if a.cfg.Media.API == "" || a.cfg.Media.URL == "" {
		return fmt.Errorf("emby host or emby API can't empty")
	}
	// 创建 emby
	emby := NewEmby(a.cfg.Media.URL, a.cfg.Media.API)
	// 上传头像
	return emby.Actor(name, acctressPath)
}

// 从javDB搜索
func (a *Actress) javDBSearch(name string) string {
	// 组合路径
	uri := fmt.Sprintf("%s/search?q=%s&f=actor", util.CheckDomainPrefix(a.JavDB), name)
	// 获取节点
	root, err := util.GetRoot(uri, a.cfg.Base.Proxy, nil)
	// 检查
	if err != nil {
		return ""
	}

	// 定义头像地址变量
	var actress string

	// 循环搜索
	root.Find(`div#actors .actor-box`).Each(func(i int, item *goquery.Selection) {
		// 获取姓名并检查
		if title, exist := item.Find(`a`).Attr("title"); exist {
			// 检查是否正确
			if strings.EqualFold(strings.TrimSpace(title), name) {
				// 获取图片并检查
				if pic, ok := item.Find("span").Attr("style"); ok {
					// 清除多余
					pic = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(pic, ")", ""), "background-image: url(", ""))
					// 检查是否空图片
					if !strings.EqualFold(pic,
						"https://javdb4.com/assets/actor_unknow-15f7d779b3d93db42c62be9460b45b79e51f8a944796eee30ed87bbb04de0a37.png") {
						actress = pic
					}

					return
				}
			}
		}
	})

	return actress
}

// 从javBus搜索
func (a *Actress) javBusSearch(name string) string {
	// 组合路径
	uri := fmt.Sprintf("%s/searchstar/%s", util.CheckDomainPrefix(a.JavBUS), name)
	// 获取节点
	root, err := util.GetRoot(uri, a.cfg.Base.Proxy, nil)
	// 检查
	if err != nil {
		return ""
	}

	// 定义头像地址变量
	var actress string

	// 循环搜索
	root.Find(`div#waterfall .item img`).Each(func(i int, item *goquery.Selection) {
		// 获取姓名并检查
		if title, exist := item.Attr("title"); exist {
			// 检查是否正确
			if strings.EqualFold(strings.TrimSpace(title), name) {
				// 获取图片并检查
				if pic, ok := item.Attr("src"); ok {
					// 检查是否空图片
					if !strings.EqualFold(pic,
						"https://pics.dmm.co.jp/mono/actjpgs/nowprinting.gif") {
						actress = pic

						return
					}
				}
			}
		}
	})

	return actress
}
