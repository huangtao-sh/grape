package text

import "strings"

// Slice 将字符串切片转换成空接口切片
func Slice(strs []string) (record []interface{}) {
	record = make([]interface{}, len(strs))
	for i, s := range strs {
		record[i] = strings.TrimSpace(s)
	}
	return
}
