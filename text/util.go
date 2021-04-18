package text

import (
	"bufio"
	"compress/gzip"
	"grape/gbk"
	"grape"
	"io"
	"os"
	"strings"
)

// Slice 将字符串切片转换成空接口切片
func Slice(strs []string) (record []interface{}) {
	record = make([]interface{}, len(strs))
	for i, s := range strs {
		record[i] = strings.TrimSpace(s)
	}
	return
}

// SplitFunc 拆分函数
type SplitFunc func(*bufio.Scanner) []string

// ConvertFunc 转换函数
type ConvertFunc func([]string) []string

// Decode 对 io.Reader 进行解码，处理 gz 文件和 GBK 编码的文件
func Decode(r io.Reader, isGz bool, isGbk bool) io.Reader {
	var err error
	if isGz {
		r, err = gzip.NewReader(r)
		grape.CheckFatal(err)
	}
	if isGbk {
		r = gbk.NewReader(r)
	}
	return r
}

// File tar、zip 压缩包获取文
type File interface {
	FileInfo() os.FileInfo
	Open() (io.ReadCloser, error)
}
