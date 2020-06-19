package util

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
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

// Wlen 统计字符串的长度，汉字按2计算
func Wlen(s string) (length int) {
	runes := []rune(s)
	length = len(runes)
	for _, r := range runes {
		if r >= 0x80 {
			length++
		}
	}
	return
}

// Sprintf format string
func Sprintf(format string, a ...interface{}) string {
	i := 0
	Pattern := regexp.MustCompile(`%.*?[sdf%]`)
	StrPattern := regexp.MustCompile(`%(-)?(\d+)?s`)
	replFunc := func(s string) (d string) {
		if s == "%%" {
			return "%"
		} else if StrPattern.MatchString(s) {
			k := StrPattern.FindStringSubmatch(s)
			d = a[i].(string)
			if k[2] != "" {
				l, _ := strconv.Atoi(k[2])
				l -= Wlen(d)
				if l > 0 {
					space := strings.Repeat(" ", l)
					if k[1] == "-" {
						d += space
					} else {
						d = space + d
					}
				}
			}
		} else {
			d = fmt.Sprintf(s, a[i])
		}
		i++
		return
	}
	return Pattern.ReplaceAllStringFunc(format, replFunc)
}
