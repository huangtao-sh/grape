package main

import (
	"archive/zip"
	"fmt"
	"grape/gbk"
	"io"
	"io/ioutil"
)

// 读取全部文件内容，并进行转码
func Read(r io.Reader, is_gbk bool) ([]byte, error) {
	if is_gbk {
		r = gbk.NewReader(r)
	}
	return ioutil.ReadAll(r)
}

// 读取文件，返回字符串
func ReadFile(r *zip.File, is_gbk bool) (s string, err error) {
	f, err := r.Open()
	if err != nil {
		return
	}
	defer f.Close()
	b, err := Read(f, is_gbk)
	if err == nil {
		s = string(b)
	}
	return
}

func main() {
	z, err := zip.OpenReader("C:/Users/huangtao/huangtao.zip")
	if err != nil {
		panic("文件不存在！")
	}
	defer z.Close()
	for i, f := range z.File {
		fmt.Println("FileNmae:", i, f.Name)
		s, _ := ReadFile(f, f.Name == "testg.txt")
		fmt.Println(s)
	}
}
