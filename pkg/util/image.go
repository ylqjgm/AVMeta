package util

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	iai "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/iai/v20180301"
)

// ConvertJPG 转换图片为jpg格式
func ConvertJPG(sourceFile, newFile string) error {
	// 打开原文件
	f, err := os.Open(sourceFile)
	// 检查错误
	if err != nil {
		return err
	}
	// 关闭连接
	defer f.Close()

	// 图片解码
	src, _, err := image.Decode(f)
	// 检查错误
	if err != nil {
		return err
	}

	// 获取图片信息
	b := src.Bounds()

	// YCBCr
	img := src.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(0, 0, b.Max.X, b.Max.Y))

	// 新建并打开新图片
	cf, err := os.OpenFile(newFile, os.O_SYNC|os.O_RDWR|os.O_CREATE, 0666)
	// 检查错误
	if err != nil {
		return err
	}
	// 关闭
	defer cf.Close()

	// 图片编码
	return jpeg.Encode(cf, img, &jpeg.Options{Quality: 100})
}

// PosterCover 图片裁剪
func PosterCover(srcPhoto, newPhoto string, cfg *ConfigStruct) error {
	// 定义各项变量
	var width, height, x int
	// 获取腾讯云人脸识别
	response, err := detectFace(srcPhoto, cfg)
	// 检查错误
	if err == nil { // 腾讯云检测到
		// 获取图片宽度
		width = int(*response.Response.ImageWidth)
		// 获取图片高度
		height = int(*response.Response.ImageHeight)
		// 获取x坐标
		x = int(*response.Response.FaceInfos[0].X)
		// 获取人脸宽度
		faceWidth := int(*response.Response.FaceInfos[0].Width)
		// 检测x坐标是否在图片右侧
		if x > width/2 {
			// x坐标设置为图片一半位置
			x = width / 2
		} else if (width/2 - x - faceWidth) < 0 { // 检测x坐标加人脸宽度是否大于图片宽度的一半
			// x坐标转换为图片一半的一半
			x = width / 2 / 2
		} else { // 默认为图片左侧
			// x坐标为0
			x = 0
		}
		// 裁剪宽度为图片一半
		width /= 2
	} else {
		// 载入图片
		img, errLoad := loadCover(srcPhoto)
		// 检查错误
		if errLoad != nil {
			return errLoad
		}

		// 获取图片边界
		b := img.Bounds()
		// 获取图片宽度并设置为一半
		width = b.Max.X / 2
		// 获取图片高度
		height = b.Max.Y
		// 将x坐标设置为0
		x = width
	}

	// 生成封面图片
	err = clipCover(srcPhoto, newPhoto, x, 0, width, height)

	return err
}

// 腾讯云免费人脸识别
func detectFace(photo string, cfg *ConfigStruct) (*iai.DetectFaceResponse, error) {
	// 图片先转换为base64
	base64, err := Base64(photo)
	// 检查错误
	if err != nil {
		return nil, err
	}

	// 初始化认证
	credential := common.NewCredential(
		cfg.Media.SecretID,
		cfg.Media.SecretKey,
	)
	// 实例化客户端配置对象
	cpf := profile.NewClientProfile()
	// 配置请求地址
	cpf.HttpProfile.Endpoint = "iai.tencentcloudapi.com"

	// 实例化人脸识别请求对象
	request := iai.NewDetectFaceRequest()
	// 请求参数, 使用base64方式请求
	params := "{\"Image\":\"" + base64 + "\"}"
	// 创建请求参数
	err = request.FromJsonString(params)
	// 检查错误
	if err != nil {
		return nil, err
	}

	// 实例化客户端对象
	client, _ := iai.NewClient(credential, "", cpf)
	// 发起请求
	response, err := client.DetectFace(request)
	// 检查是否为API错误
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		err = fmt.Errorf("an API error has returned: %s", err)
	}

	return response, err
}

// 剪切图片
func clipCover(srcFile, newFile string, x, y, w, h int) error {
	// 载入图片
	src, err := loadCover(srcFile)
	// 检查错误
	if err != nil {
		return err
	}

	// 剪切图片
	img := src.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(x, y, x+w, y+h))

	// 保存图片
	saveErr := saveCover(newFile, img)
	// 检查错误
	if saveErr != nil {
		return err
	}

	return nil
}

// 载入图片
func loadCover(photo string) (img image.Image, err error) {
	// 打开图片文件
	file, err := os.Open(photo)
	// 检查错误
	if err != nil {
		return
	}
	// 关闭
	defer file.Close()

	// 图片解码
	img, _, err = image.Decode(file)

	return
}

// 保存图片
func saveCover(path string, img image.Image) error {
	// 新建并打开文件
	f, err := os.OpenFile(path, os.O_SYNC|os.O_RDWR|os.O_CREATE, 0666)
	// 检查错误
	if err != nil {
		return err
	}
	// 关闭
	defer f.Close()

	// 获取文件后缀
	ext := filepath.Ext(path)

	// 如果是jpeg类型
	if strings.EqualFold(ext, ".jpg") || strings.EqualFold(ext, ".jpeg") {
		// jpeg图片编码
		err = jpeg.Encode(f, img, &jpeg.Options{Quality: 80})
	} else if strings.EqualFold(ext, ".png") { // png类型
		// png图片编码
		err = png.Encode(f, img)
	}

	return err
}
