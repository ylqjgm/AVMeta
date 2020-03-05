package actress

import (
	"fmt"
	"strings"

	"github.com/ylqjgm/AVMeta/pkg/util"

	"github.com/PuerkitoBio/goquery"
)

var (
	javBusCensored   = "actresses/%d"
	javBusUnCensored = "uncensored/actresses/%d"
)

// JavBUS 采集
func JavBUS(site, proxy string, page int, censored bool) (actress map[string]string, nextPage bool, err error) {
	// 定义地址变量
	var uri string
	// 判断是有码还是无码
	if censored {
		// 有码女优列表地址
		uri = fmt.Sprintf("%s/%s", util.CheckDomainPrefix(site), fmt.Sprintf(javBusCensored, page))
	} else {
		uri = fmt.Sprintf("%s/%s", util.CheckDomainPrefix(site), fmt.Sprintf(javBusUnCensored, page))
	}

	// 定义女优列表
	actress = make(map[string]string)

	// 打开女优列表
	root, err := util.GetRoot(uri, proxy, nil)
	// 检查错误
	if err != nil {
		return nil, false, err
	}

	// 获取
	root.Find(`.item a`).Each(func(i int, item *goquery.Selection) {
		// 获取名字
		name := strings.TrimSpace(item.Find(`.photo-info span`).Text())
		// 获取头像
		face, _ := item.Find(`.photo-frame img`).Attr("src")
		// 清除多余
		face = strings.TrimSpace(face)
		// 是否获取到
		if name != "" && face != "" {
			actress[name] = face
		}
	})

	// 查询下一页节点
	_, exists := root.Find(`a#next`).Attr("href")
	// 找到
	if exists {
		return actress, true, nil
	}

	return actress, false, nil
}
