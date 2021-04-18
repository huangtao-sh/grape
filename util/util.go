package util

import (
	"fmt"
	"grape"
)

// Dater 数据接口
type Dater interface {
	Next() bool
	Read() []interface{}
}

// Println 打印一行数据
func Println(data Dater) {
	var row []interface{}
	for data.Next() {
		row = data.Read()
		fmt.Println(row...)
	}
}

// Printf 格式打印
func Printf(format string, data Dater) {
	var row []interface{}
	for data.Next() {
		row = data.Read()
		fmt.Print(grape.Sprintf(format, row...))
	}
}
