package util

import (
	"path"
	"path/filepath"
	"strings"
)

// GetCode 提取番号
func GetCode(filename, filter string) string {
	// 获取正确文件名
	filename = filepath.Base(strings.ToLower(filename))
	// 删除扩展名
	filename = strings.TrimSuffix(filename, path.Ext(filename))
	// 转换过滤规则为数组
	filters := strings.Split(filter, "||")
	// 循环过滤
	for _, f := range filters {
		// 过滤
		filename = strings.ReplaceAll(filename, f, "")
	}
	// 将所有 . 替换为 -
	filename = strings.ReplaceAll(filename, ".", "-")
	// 过滤空格
	filename = strings.TrimSpace(filename)

	return filename
}

// GetNumberPath 获取正确的保存路径
func GetNumberPath(replaceStr map[string]string, cfg *ConfigStruct) string {
	// 获取运行路径
	base := GetRunPath()
	// 组合路径
	base = base + "/" + cfg.Path.Success
	// 获取保存规则
	rule := cfg.Path.Directory
	// 循环替换
	for key, val := range replaceStr {
		rule = strings.ReplaceAll(rule, key, val)
	}

	// 定义特殊字符数组
	filter := []string{"\\", ":", "*", "?", `"`, "<", ">", "|"}
	// 循环过滤
	for _, v := range filter {
		rule = strings.ReplaceAll(rule, v, "")
	}
	// 多余的反斜线
	rule = strings.ReplaceAll(rule, "//", "/")

	return base + "/" + rule
}

// CheckDomainPrefix 检查域名最后的斜线
func CheckDomainPrefix(domain string) string {
	// 是否为空
	if domain == "" {
		return ""
	}

	// 获取最后一个字符
	last := domain[len(domain)-1:]
	// 如果是斜线
	if last == "/" {
		domain = domain[:len(domain)-1]
	}

	return domain
}

// IntroFilter 过滤简介
func IntroFilter(intro string) string {
	// 替换<br>
	intro = strings.ReplaceAll(intro, "<br>", "\n")
	intro = strings.ReplaceAll(intro, "<br/>", "\n")
	intro = strings.ReplaceAll(intro, "<br />", "\n")
	// 替换\r\n
	intro = strings.ReplaceAll(intro, "\r\n", "\n")
	// 替换\r
	intro = strings.ReplaceAll(intro, "\r", "\n")
	// 替换\n\n
	intro = strings.ReplaceAll(intro, "\n\n", "\n")

	// 清除多余空白
	return strings.TrimSpace(intro)
}
