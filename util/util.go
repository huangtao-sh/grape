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

// Recover 检查错误，并打印
func Recover() {
	if r := recover(); r != nil {
		fmt.Println("错误：", r)
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
		fmt.Print(Sprintf(format, row...))
	}
}

// Wlen 统计字符串的打印宽度， ASCII 为 1，汉字为 2
func Wlen(s string) (length int) {
	runes := []rune(s)  // 先转换成 Unicode 字
	length = len(runes) // 统计长度
	for _, r := range runes {
		if r >= 0x80 {
			length++ // 每个汉字再额外加 1
		}
	}
	return
}

// Reverse 对 byte Slice 进行反转
func Reverse(s []byte) []byte {
	for i, l := 0, len(s)-1; i < l; {
		s[i], s[l] = s[l], s[i]
		l--
		i++
	}
	return s
}

// formatInt 对整数进行格式化，增加千分节
func formatInt(k string) string {
	var d []byte
	s := []byte(strings.TrimSpace(k))
	flag := false // 是否为负数
	if s[0] == '-' {
		s = s[1:]
		flag = true
	}
	s = Reverse(s)
	for i, c := range s {
		if i%3 == 0 && i > 0 {
			d = append(d, ',')
		}
		d = append(d, c)
	}
	if flag {
		d = append(d, '-')
	}
	return string(Reverse(d))
}

// Sprintf 字符串格式化，解决汉字宽度及数字无千分节问题
func Sprintf(format string, args ...interface{}) (d string) {
	StrPattern := regexp.MustCompile(`%(-)?(\d+)s`)
	IntPattern := regexp.MustCompile(`%(\d+),d`)
	FloatPattern := regexp.MustCompile(`%(\d+),\.(\d+)f`)
	ValuesPattern := regexp.MustCompile(`%(\d+)V`)
	Pattern := regexp.MustCompile(`%.*?[sdfvV%]`)
	i := 0
	replFunc := func(s string) (d string) {
		if s == "%%" {
			return "%"
		} else if k := StrPattern.FindStringSubmatch(s); k != nil {
			d = args[i].(string)
			l, _ := strconv.Atoi(k[2])
			if l-Wlen(d) > 0 {
				space := strings.Repeat(" ", l-Wlen(d))
				if k[1] == "-" {
					d = d + space
				} else {
					d = space + d
				}
			}
		} else if k := IntPattern.FindStringSubmatch(s); k != nil {
			l, _ := strconv.Atoi(k[1])
			d = string(formatInt(fmt.Sprintf("%d", args[i])))
			space := strings.Repeat(" ", l-len(d))
			d = space + d
		} else if k := FloatPattern.FindStringSubmatch(s); k != nil {
			l, _ := strconv.Atoi(k[1])
			s, _ := strconv.Atoi(k[2])
			d = fmt.Sprintf(fmt.Sprintf("%%.%df", s), args[i])
			a := strings.Split(d, ".")
			a[0] = formatInt(a[0])
			d = strings.Join(a, ".")
			space := strings.Repeat(" ", l-len(d))
			d = space + d
		} else if k := ValuesPattern.FindStringSubmatch(s); k != nil {
			count, _ := strconv.Atoi(k[1])
			d = strings.Repeat("?,", count-1)
			d += "?"
			return fmt.Sprintf("Values(%s)", d)
		} else {
			d = fmt.Sprintf(s, args[i])
		}
		i++
		return
	}
	d = Pattern.ReplaceAllStringFunc(format, replFunc)
	return
}

// Match 匹配字符串
func Match(pattern string, s string) (matched bool) {
	matched, err := regexp.MatchString(pattern, s)
	CheckFatal(err)
	return
}

// FullMatch 全匹配字符串
func FullMatch(pattern string, s string) bool {
	return Match(fmt.Sprintf("^%s$", pattern), s)
}
