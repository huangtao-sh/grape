package text

import (
	"bufio"
	"io"
)

// BasicReader 读取文本的基础实例
type BasicReader struct {
	bufio.Scanner
	record []interface{}
}

// NewTestReader 创建 Reader
func NewTestReader(r io.Reader) BasicReader {
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
	return &FixedReader{NewTestReader(r), offsets}
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

// Reader 读取数据模块
type Reader struct {
	*bufio.Scanner
	Split     SplitFunc
	converter []ConvertFunc
}

// NewReader 数据模块构造函数
func NewReader(r io.Reader, skipHeader bool, split SplitFunc, converters ...ConvertFunc) *Reader {
	scanner := bufio.NewScanner(r)
	if skipHeader {
		scanner.Scan()
	}
	return &Reader{scanner, split, converters}
}

// Next 是否有下一条数据
func (r *Reader) Next() bool {
	return r.Scan()
}

// Read 读取当前数据
func (r *Reader) Read() []interface{} {
	var row []string
	row = r.Split(r.Scanner)
	for _, convert := range r.converter {
		row = convert(row)
		if row == nil {
			return nil
		}
	}
	return Slice(row)
}
