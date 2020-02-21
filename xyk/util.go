package xyk

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Atoi 字符串转换成整数
func Atoi(s string) (result int) {
	s = strings.Replace(s, ".", "", 1)
	s = strings.Replace(s, "+", "", 1)
	result, err := strconv.Atoi(s)
	if err != nil {
		panic("转换数字失败")
	}
	return
}

// DateAdd 增加日期
func DateAdd(date string, days int) (result string) {
	d, _ := time.Parse("2006-01-02 15:04:05", date[0:4]+"-"+date[4:6]+"-"+date[6:8]+" 00:00:00")
	d = d.AddDate(0, 0, days)
	return fmt.Sprintf("%04d%02d%02d", d.Year(), d.Month(), d.Day())
}

// PrevDay 上一天
func PrevDay(date string) string {
	return DateAdd(date, -1)
}

// NextDay 下一天
func NextDay(date string) string {
	return DateAdd(date, 1)
}

// SplitData 拆分数据
func SplitData(bytes []byte, offsets []int, columns []int) (result []string) {
	for _, index := range columns {
		start := offsets[index]
		end := offsets[index+1] 
		result = append(result, strings.TrimSpace(string(bytes[start:end])))
	}
	return
}
