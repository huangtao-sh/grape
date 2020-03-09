package util

import (
	"fmt"
	"log"
	"os"
)

// CheckErr 检查是否有错误，并退出操作系统
func CheckErr(err error, exitCode int) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(exitCode)
	}
}

// CheckFatal 检查致命错误
func CheckFatal(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
