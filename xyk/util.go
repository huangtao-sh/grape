package xyk

import (
	"strconv"
	"strings"
)

// 把小数转换成整数，直接移除小数点
func Atoi(s string) (result int) {
	s = strings.Replace(s, ".", "", 1)
	s = strings.Replace(s, "+", "", 1)
	result, err := strconv.Atoi(s)
	if err != nil {
		panic("转换数字失败")
	}
	return
}
