package date

import (
	"fmt"
	"regexp"
	"time"
)

// Weekdays 中文星期的表示
var Weekdays [7]string = [7]string{
	"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"}

// Date Basic type
type Date struct {
	time.Time
}

// Today 返回当前日期
func Today() Date {
	return Date{time.Now()}
}

// NewDate 构造日期
func NewDate(s string) (date Date) {
	var parts [3]string
	if regexp.MustCompile(`\d{4}\d{2}\d{2}`).MatchString(s) {
		parts[0] = s[:4]
		parts[1] = s[4:6]
		parts[2] = s[6:]
	} else if regexp.MustCompile(`\d{4}-\d{1,2}-\d{1,2}`).MatchString(s) {

	}
	return
}

// Format 格式化日期
// %%  %
// %Y  4位年份，如：2020
// %M  2位月份，如：02
// %D  2位日期，如：15
// %Q  1位季度，如：2
// %W  星期，如星期一、星期二、星期日...
func (d Date) Format(format string) (s string) {
	f := func(s string) (r string) {
		switch s {
		case "%%":
			r = "%"
		case "%Y":
			r = fmt.Sprintf("%04d", d.Year())
		case "%y":
			r = fmt.Sprintf("%02d", d.Year()%100)
		case "%M":
			r = fmt.Sprintf("%02d", d.Month())
		case "%D":
			r = fmt.Sprintf("%02d", d.Day())
		case "%Q":
			r = fmt.Sprintf("%1d", d.Quater())
		case "%F":
			r = fmt.Sprintf("%04d-%02d-%02d", d.Year(), d.Month(), d.Day())
		case "%W":
			r = Weekdays[d.Weekday()]
		default:
			panic("Wrong Format")
		}
		return
	}
	r := regexp.MustCompile("%.")
	s = r.ReplaceAllStringFunc(format, f)
	return
}

// Add 增加
func (d Date) Add(years, months, days int) (date Date) {
	date = Date{d.AddDate(years, months, days)}
	return
}

// Quater 季度
func (d Date) Quater() int {
	return (int(d.Month()) + 2) / 3
}

// String 返回字符串
func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year(), d.Month(), d.Day())
}

// NextDay 下一日
func (d Date) NextDay() Date {
	return d.Add(0, 0, 1)
}

// PrevDay 上一日
func (d Date) PrevDay() Date {
	return d.Add(0, 0, -1)
}
