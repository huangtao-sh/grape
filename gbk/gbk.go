// GBK 转码包
package gbk

import (
	"io"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var (
	Encoder *encoding.Encoder
	Decoder *encoding.Decoder
)

func init() {
	Encoder = simplifiedchinese.GBK.NewEncoder()
	Decoder = simplifiedchinese.GBK.NewDecoder()
}

// 新建一个 GBK 编码的 Reader
func NewReader(r io.Reader) *transform.Reader {
	return transform.NewReader(r, Decoder)
}

// 新建一个 GBK 编码的 Writer
func NewWriter(w io.Writer) *transform.Writer {
	return transform.NewWriter(w, Encoder)
}

// 将 UTF8 字符串转换成 GBK 编码的 bytes
func Encode(s string) ([]byte, error) {
	return Encoder.Bytes([]byte(s))
}

//将 GBK 编码的 bytes 转换成 UTF8 字符串
func Decode(b []byte) (string, error) {
	return Decoder.String(string(b))
}
