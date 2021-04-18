// Package gbk  GBK转码包
// Writer : Huang Tao 2020/03/01
// 支持将 GBK 编码的二进制解码为字符串，或将字符串使用 GBK 编码成二进制
package grape

import (
	"io"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var (
	// Encoder GBK 编码器
	GBKEncoder *encoding.Encoder
	// Decoder GBK 解码器
	GBKDecoder *encoding.Decoder
)

func init() {
	GBKEncoder = simplifiedchinese.GBK.NewEncoder()
	GBKDecoder = simplifiedchinese.GBK.NewDecoder()
}

// NewReader 新建一个 GBK 编码的 Reader
func NewGBKReader(r io.Reader) *transform.Reader {
	return transform.NewReader(r, GBKDecoder)
}

// NewGBKWriter 新建一个 GBK 编码的 Writer
func NewGBKWriter(w io.Writer) *transform.Writer {
	return transform.NewWriter(w, GBKEncoder)
}

// GBKEncode 将 UTF8 字符串转换成 GBK 编码的 bytes
func GBKEncode(s string) ([]byte, error) {
	return GBKEncoder.Bytes([]byte(s))
}

// GBKDecode 将 GBK 编码的 bytes 转换成 UTF8 字符串
func GBKDecode(b []byte) (string, error) {
	return GBKDecoder.String(string(b))
}
