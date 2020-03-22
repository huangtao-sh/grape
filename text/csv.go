package text

import (
	"encoding/csv"
	"io"
)

// CsvReader 读取 Csv 文件
type CsvReader struct {
	csv.Reader
	record []interface{}
}

// NewCsvReader 新建CsvReader
func NewCsvReader(r io.Reader) *CsvReader {
	return &CsvReader{Reader: *csv.NewReader(r)}
}

// Next 判断是否还有数据
func (r *CsvReader) Next() (exists bool) {
	rec, err := r.Reader.Read()
	if err != nil {
		r.record = Slice(rec)
		exists = true
	} else if err == io.EOF {
		exists = false
	} else {
		panic(err.Error())
	}
	return
}

// Read 读取一行记录
func (r *CsvReader) Read() []interface{} {
	return r.record
}
