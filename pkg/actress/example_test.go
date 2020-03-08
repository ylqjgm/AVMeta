package actress_test

import (
	"github.com/ylqjgm/AVMeta/pkg/actress"
)

// 下载女优头像，第一个参数为下载网站，
// 第二个参数为开始页数，第三个参数为是否有码女优
func ExampleActress_Fetch() {
	// 通过该方法下载的女优头像，均自动保存到程序执行目录
	// 下的 actress 文件夹中

	// 初始化一个Actress对象
	act := actress.NewActress()
	// 下载 javbus 网站有码女优，从第 1 页开始
	act.Fetch(actress.JAVBUS, 1, true)
	// 下载 javbus 网站无码女优，从第 2 页开始
	act.Fetch(actress.JAVBUS, 2, false)
	// 下载 javdb 网站有码女优，从第 10 页开始
	act.Fetch(actress.JAVDB, 10, true)
	// 下载 javdb 网站无码女优，从第 20 页开始
	act.Fetch(actress.JAVDB, 20, false)
}

// 将本地女优头像入库到 Emby中，
// 女优头像需存储在程序执行目录的
// actress 目录中，且以 "女优名字.jpg" 格式保存
func ExampleActress_Put() {
	// 初始化一个Actress对象
	act := actress.NewActress()
	// 入库头像
	act.Put()
}

// 单个女优头像入库，参数一为女优名称，与 Emby 中对应，
// 参数二为女优头像保存路径
func ExampleEmby_Actor() {
	// 初始化一个Emby对象，传入两个字符串，
	// hostURL 为 Emby 访问地址，
	// apiKey 为 获取的 API Key
	emby := actress.NewEmby("http://127.0.0.1:8096", "123456")
	// 入库头像
	err := emby.Actor("上原亜衣", "./actress/上原亜衣.jpg")
	if err != nil {
		panic(err)
	}
}
