package util

import (
	"fmt"
	"os"
)

// CheckErr 检查是否有错误，并退出操作系统
func CheckErr(err error, exitCode int) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(exitCode)
	}
}

// CheckPanic 检查致命错误
func CheckPanic(err error) {
	if err != nil {
		panic(err.Error())
	}
}

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
		fmt.Printf(format, row...)
	}
}
