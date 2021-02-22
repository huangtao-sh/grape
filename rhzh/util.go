package rhzh

import (
	"strings"
)

// FullChar 将含有半角的字符串转换为全角
func FullChar(s string) string {
	r := []rune(s)
	for i, v := range r {
		if v > 0x20 && v < 0x7F {
			r[i] = v + 0xFF00 - 0x20
		}
	}
	return string(r)
}

// Date 转换日期，将格式为 YYYYMMDD 的日期转换成 YYYY-MM-DD
func Date(date string) string {
	if len(date) == 8 {
		date = strings.Join([]string{date[:4], date[4:6], date[6:]}, "-")
	} else {
		date = ""
	}
	return date
}
