package csvfile

import (
	"encoding/csv"
	"io"
	"os"

	"grape/gbk"
)

// CSV Reader
// 逐行读取 CSV 文件，每次读取一条记录
type CsvReader struct {
	f *os.File
	r *csv.Reader
}

// 关闭 CSV Reader
func (r CsvReader) Close() error {
	return r.f.Close()
}

// 将字符串切片转换成空接口切片
func Conv(strs []string) []interface{} {
	result := make([]interface{}, len(strs))
	for i, s := range strs {
		result[i] = s
	}
	return result
}

// 从 CSV 文件中读取一行数据，成并转换成空接口切片
func (r CsvReader) Read() (record []interface{}, err error) {
	rec, err := r.r.Read()
	if err == nil {
		record = Conv(rec)
	}
	return
}

// 创建一个 CSV Reader
func NewReader(filename string, encoding string) (reader *CsvReader, err error) {
	var r io.Reader
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	if encoding == "GBK" {
		r = gbk.Reader(f)
	} else {
		r = f
	}
	reader = &CsvReader{f, csv.NewReader(r)}
	return
}
