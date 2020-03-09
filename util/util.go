package util

import (
	"fmt"
	"os"
)

// CheckErr 检查是否有错误，并退出操作系统
func CheckErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
