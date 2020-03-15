package text

import (
	"bufio"
	"io"
)

// Slice 将字符串切片转换成空接口切片
func Slice(strs []string) []interface{} {
	result := make([]interface{}, len(strs))
	for i, s := range strs {
		result[i] = s
	}
	return result
}

// BasicReader 读取文本的基础实例
type BasicReader struct {
	bufio.Scanner
	record []interface{}
}

// NewReader 创建 Reader
func NewReader(r io.Reader) BasicReader {
	scanner := bufio.NewScanner(r)
	return BasicReader{Scanner: *scanner}
}

// SetRecord 设置当前记录的值
func (r *BasicReader) SetRecord(record []interface{}) {
	r.record = record
}

// Next 获取下一条数据
func (r *BasicReader) Next() bool {
	return r.record != nil
}

// Read 读取当前数据
func (r *BasicReader) Read() []interface{} {
	return r.record
}

// FixedReader 固定偏移量的文本读取
type FixedReader struct {
	BasicReader
	offsets []int
}

// NewFixedReader 创建固定偏移量的文本读取
func NewFixedReader(r io.Reader, offsets []int) *FixedReader {
	return &FixedReader{NewReader(r), offsets}
}

// Next 读取下一行数据
func (r *FixedReader) Next() bool {
	var record []interface{}
	if r.Scan() {
		row := r.Bytes()
		for i := range r.offsets[1:] {
			record = append(record, string(row[r.offsets[i]:r.offsets[i+1]]))
		}
		r.record = record
		return true
	}
	return false
}
