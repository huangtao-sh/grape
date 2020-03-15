package text

import (
	"encoding/csv"
	"grape/util"
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
func (r *CsvReader) Next() bool {
	rec, err := r.Reader.Read()
	if err == io.EOF {
		return false
	}
	util.CheckFatal(err)
	r.record = Slice(rec)
	return true
}

// Read 读取一行记录
func (r *CsvReader) Read() []interface{} {
	return r.record
}
