package util_test

import (
	"fmt"

	"github.com/ylqjgm/AVMeta/pkg/util"
)

// 读取配置文件
func ExampleGetConfig() {
	cfg, err := util.GetConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg.Base.Success)

	// Output:
	// success
}

// 写入默认配置文件
func ExampleWriteConfig() {
	cfg, err := util.WriteConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg.Base.Success)

	// Output:
	// success
}

// 遍历目录
func ExampleWalkDir() {
	files, err := util.WalkDir("./", "success", "fail")
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		fmt.Println(f)
	}

	// Output:
	// ./xxx
	// .....
}

// 获取程序执行路径
func ExampleGetRunPath() {
	cur := util.GetRunPath()
	fmt.Println(cur)

	// Output:
	// AVMeta run path
}

// Base64 编码
func ExampleBase64() {
	enc, err := util.Base64("./AVMeta.go")
	if err != nil {
		panic(err)
	}
	fmt.Println(enc)

	// Output:
	// AVMeta.go base64 string...
}

// 失败文件操作
func ExampleFailFile() {
	util.FailFile("./a.mp4", "fail")
}

// 移动文件
func ExampleMoveFile() {
	err := util.MoveFile("./a.mp4", "./b.mp4")
	if err != nil {
		panic(err)
	}
}

// 获取文件大小
func ExampleGetFileSize() {
	size := util.GetFileSize("./AVMeta.go")
	fmt.Println(size)
}

// 写入文件
func ExampleWriteFile() {
	err := util.WriteFile("./a.txt", []byte("aaa"))
	if err != nil {
		panic(err)
	}
}

// 文件是否存在
func ExampleExists() {
	fmt.Println(util.Exists("./a.txt"))
}

// 创建请求
func ExampleMakeRequest() {
	data, status, err := util.MakeRequest("GET", "https://www.baidu.com", "", nil, nil, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(status)
	fmt.Println(string(data))
}

// 获取远程数据
func ExampleGetResult() {
	data, err := util.GetResult("https://www.baidu.com", "", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}

// 获取远程树结构
func ExampleGetRoot() {
	root, err := util.GetRoot("https://www.baidu.com", "", nil)
	if err != nil {
		panic(err)
	}

	html, err := root.Html()
	if err != nil {
		panic(err)
	}

	fmt.Println(html)
}

// 下载远程图片
func ExampleSavePhoto() {
	err := util.SavePhoto("https://pics.javbus.com/actress/okq_a.jpg", "./三上悠亜.jpg", "", false)
	if err != nil {
		panic(err)
	}
}

// 转换图片
func ExampleConvertJPG() {
	err := util.ConvertJPG("./a.png", "./a.jpg")
	if err != nil {
		panic(err)
	}
}

// 裁剪图片
func ExamplePosterCover() {
	cfg, err := util.GetConfig()
	if err != nil {
		panic(err)
	}

	err = util.PosterCover("./fanart.jpg", "./poster.jpg", cfg)
	if err != nil {
		panic(err)
	}
}

// 提取番号信息
func ExampleGetCode() {
	number := util.GetCode("./BF-592_xyz.mp4", "_||xyz")
	fmt.Println(number)

	// Output:
	// BF-592
}

// 获取准确保存路径
func ExampleGetNumberPath() {
	cfg := &util.ConfigStruct{
		Path: PathStruct{
			Success:   "success",
			Fail:      "fail",
			Directory: "{actor}/{number}",
		},
	}
	filter := make(map[string]string)
	filter["actor"] = "三上悠亜"
	filter["number"] = "SSNI-703"

	fmt.Println(util.GetNumberPath(filter, cfg))

	// Output:
	// ./三上悠亜/SSNI-703
}

// 检查域名斜线
func ExampleCheckDomainPrefix() {
	fmt.Println(util.CheckDomainPrefix("https://www.baidu.com/"))

	// Output:
	// https://www.baidu.com
}

// 简介过滤
func ExampleIntroFilter() {
	intro := "第一段<br>第二段<br />第三段"
	fmt.Println(util.IntroFilter(intro))

	// Output:
	// 第一段
	// 第二段
	// 第三段
}
