package data

import "grape/text"

// ConvertFunc 转换函数
type ConvertFunc func([]string) []string

// Reader 数据读取接口
type Reader interface {
	Next() bool
	Read() []string
}

// ConvertReader 带通道的数据读取
type ConvertReader struct {
	Reader
	converters []ConvertFunc
}

// NewConvertReader  DReader 构造函数
func NewConvertReader(r Reader, converters ...ConvertFunc) *ConvertReader {
	return &ConvertReader{r, converters}
}

// ReadAll 读取所有数据
func (r *ConvertReader) ReadAll(d text.Data) {
	defer d.Close()
	var row []string
	for r.Next() {
		row = r.Read()
		for _, convert := range r.converters {
			row = convert(row)
			if row == nil {
				break
			}
		}
		if row != nil {
			d.Write(text.Slice(row)...)
		}
	}
}
