package text

import (
	"strconv"
)

// Include ConvertFunc 包含指定的列
func Include(columns ...int) ConvertFunc {
	return func(s []string) (d []string) {
		d = make([]string, len(columns))
		for i, idx := range columns {
			d[i] = s[idx]
		}
		return
	}
}

// Exclude ConvertFunc 排除指定的列
func Exclude(columns ...int) ConvertFunc {
	Columns := make(map[int]bool)
	for _, i := range columns {
		Columns[i] = true
	}
	return func(s []string) (d []string) {
		for i, value := range s {
			if !Columns[i] {
				d = append(d, value)
			}
		}
		return
	}
}

// UnQuote ConvertFunc 删除字符串的引号
func UnQuote(s []string) (d []string) {
	var err error
	d = make([]string, len(s))
	for i := range s {
		d[i], err = strconv.Unquote(s[i])
		if err != nil {
			d[i] = s[i]
		}

	}
	return
}
