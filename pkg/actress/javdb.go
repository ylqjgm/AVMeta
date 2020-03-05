package actress

import (
	"fmt"
	"strings"

	"github.com/ylqjgm/AVMeta/pkg/util"

	"github.com/PuerkitoBio/goquery"
)

var (
	javDBCensored   = "actors?page=%d"
	javDBUnCensored = "actors/uncensored?page=%d"
)

// JavDB 采集
func JavDB(site, proxy string, page int, censored bool) (actress map[string]string, nextPage bool, err error) {
	// 定义地址变量
	var uri string
	// 判断是有码还是无码
	if censored {
		// 有码女优列表地址
		uri = fmt.Sprintf("%s/%s", util.CheckDomainPrefix(site), fmt.Sprintf(javDBCensored, page))
	} else {
		uri = fmt.Sprintf("%s/%s", util.CheckDomainPrefix(site), fmt.Sprintf(javDBUnCensored, page))
	}

	// 定义女优列表
	actress = make(map[string]string)

	// 打开女优列表
	root, err := util.GetRoot(uri, proxy, nil)
	// 检查错误
	if err != nil {
		return nil, false, err
	}

	// 获取搜索
	root.Find(`.actor-box a`).Each(func(i int, item *goquery.Selection) {
		// 获取姓名
		name := strings.TrimSpace(item.Find(`strong`).Text())
		// 获取头像
		face, _ := item.Find(`.image span`).Attr("style")
		// 清除多余
		face = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(face, ")", ""), "background-image: url(", ""))
		// 是否获取到
		if name != "" && face != "" {
			// 加入列表
			actress[name] = face
		}
	})

	// 查询下一页节点
	_, exists := root.Find(`a.pagination-next`).Attr("href")
	// 找到
	if exists {
		return actress, true, nil
	}

	return actress, false, nil
}
