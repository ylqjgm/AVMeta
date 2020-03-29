package media

import (
	"bytes"
	"encoding/binary"
	"math"
	"strconv"

	"github.com/ylqjgm/AVMeta/pkg/util"
)

// VSMeta 群晖元数据结构体
type VSMeta struct {
	B bytes.Buffer
}

// NewVSMeta 实例化一个VSMeta
func NewVSMeta() *VSMeta {
	return &VSMeta{B: bytes.Buffer{}}
}

// ParseVSMeta 将刮削对象转换为 VSMeta 所需字节集
//
// s IScraper刮削接口，传入刮削对象
func (v *VSMeta) ParseVSMeta(m *Media) {
	// 写入头部
	v.writeTag(0x08)
	// 写入文件类型
	v.writeTag(0x01)
	// 写入标题1
	v.writeTag(0x12)
	v.writeString(v.subStr(m.Title.Inner, 255, 85))
	// 写入标题2
	v.writeTag(0x1A)
	v.writeString(v.subStr(m.Title.Inner, 255, 85))
	// 写入副标题
	if m.SortTitle != "" {
		v.writeTag(0x22)
		v.writeString(v.subStr(m.SortTitle, 255, 85))
	}

	// 写入年份
	if m.Year != "" {
		year, _ := strconv.Atoi(m.Year)
		v.writeTag(0x28)
		v.writeInt(year)
		v.writeTag(0x0F)
	}

	// 写入日期
	if m.Premiered != "" {
		v.writeTag(0x32)
		v.writeString(m.Premiered)
	}
	// 写入锁定
	v.writeTag(0x38)
	v.writeTag(0x01)
	// 写入简介
	if m.Plot.Inner != "" {
		v.writeTag(0x42)
		v.writeString(v.subStr(m.Plot.Inner, 4096, 1360))
	}
	// 写入来源json
	v.writeTag(0x4A)
	v.writeString("null")

	// 写入组数据
	v.writeGroup(m)

	// 写入级别
	v.writeTag(0x5A)
	v.writeString("XXX")
	// 写入评分
	v.writeTag(0x60)
	v.writeInt(0)
}

// 写入标签
func (v *VSMeta) writeTag(tag uint8) {
	binary.Write(&v.B, binary.LittleEndian, tag)
}

// 写入数字
func (v *VSMeta) writeInt(num int) {
	binary.Write(&v.B, binary.LittleEndian, uint8(num))
}

// 写入字符串
func (v *VSMeta) writeString(str string) {
	v.writeBytes([]byte(str))
}

// 写入字节集
func (v *VSMeta) writeBytes(bs []byte) {
	// 先写入长度
	v.writeLength(len(bs))
	// 写入字节集
	binary.Write(&v.B, binary.LittleEndian, bs)
}

// 写入长度
//
// length整型，要写入的长度
func (v *VSMeta) writeLength(length int) {
	v.writeLengthByBuffer(length, &v.B)
}

// 写入长度到指定buffer
func (v *VSMeta) writeLengthByBuffer(length int, buf *bytes.Buffer) {
	// 只要长度大于128则循环处理
	for length > 128 {
		// 取余
		binary.Write(buf, binary.LittleEndian, byte((length%128)+128))
		// 长度处理
		length = int(math.Floor(float64(length) / float64(128)))
	}
	// 写入长度
	binary.Write(buf, binary.LittleEndian, byte(length))
}

// 写入组数据
func (v *VSMeta) writeGroup(m *Media) {
	// 定义字节集
	var buf []byte
	// 是否有演员
	if len(m.Actor) > 0 {
		// 循环演员信息
		for _, val := range m.Actor {
			// 写入演员信息
			buf = append(buf, uint8(0x0A), byte(len(val.Name)))
			buf = append(buf, []byte(val.Name)...)
		}
	}
	// 是否有导演
	if m.Director.Inner != "" {
		// 写入导演信息
		buf = append(buf, uint8(0x12), byte(len(m.Director.Inner)))
		buf = append(buf, []byte(m.Director.Inner)...)
	}
	// 是否有类型
	if len(m.Genre) > 0 {
		// 循环类型信息
		for _, val := range m.Genre {
			// 写入类型信息
			buf = append(buf, uint8(0x1A), byte(len(val.Inner)))
			buf = append(buf, []byte(val.Inner)...)
		}
	}
	// 是否有厂商
	if m.Studio.Inner != "" {
		// 写入编剧信息
		buf = append(buf, uint8(0x22), byte(len(m.Studio.Inner)))
		buf = append(buf, []byte(m.Studio.Inner)...)
	}
	// 是否有数据
	if len(buf) > 0 {
		// 写入标签
		v.writeTag(0x52)
		// 写入数据
		v.writeBytes(buf)
	}
}

// 写入封面
//
// file 字符串，封面图片路径
func (v *VSMeta) writePoster(file string) {
	// 获取封面base64
	poster, err := util.Base64(file)
	// 检查错误
	if err != nil {
		return
	}
	// 获取封面md5
	has := util.MD5String(poster)

	// 写入封面
	v.writeTag(0x8A)
	v.writeTag(0x01)
	v.writeString(poster)

	// 写入md5
	v.writeTag(0x92)
	v.writeTag(0x01)
	v.writeString(has)
}

// 写入背景
//
// file 字符串，背景图片路径
func (v *VSMeta) writeFanart(file string) {
	// 获取背景base64
	fanart, err := util.Base64(file)
	// 检查错误
	if err != nil {
		return
	}
	// 获取背景md5
	has := util.MD5String(fanart)

	// 定义字节集
	var buf []byte
	length := len(fanart)
	// 写入标签
	buf = append(buf, uint8(0x0A))
	// 写入长度
	for length > 128 {
		buf = append(buf, byte((length%128)+128))
		length = int(math.Floor(float64(length) / float64(128)))
	}
	buf = append(buf, byte(length))
	// 写入内容
	buf = append(buf, []byte(fanart)...)
	// 写入标签
	buf = append(buf, uint8(0x12), uint8(0x20))
	// 写入内容
	buf = append(buf, []byte(has)...)

	// 写入组2
	v.writeTag(0xAA)
	// 写入标签
	v.writeTag(0x01)
	// 写入内容
	v.writeBytes(buf)
}

// 截取字符串
func (v *VSMeta) subStr(str string, max, length int) string {
	if len(str) > max {
		return string([]rune(str)[:length])
	}

	return str
}
