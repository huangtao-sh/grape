package text

import (
	"bufio"
	"strings"
)

// OffsetSpliter 定位分割
type OffsetSpliter struct {
	offsets []int
}

// NewOffsetSpliter OffsetSpliter 构造函数
func NewOffsetSpliter(offsets []int) *OffsetSpliter {
	return &OffsetSpliter{offsets}
}

// Split 执行拆分
func (s *OffsetSpliter) Split(bytes []byte) (record []interface{}) {
	offsets := s.offsets
	var i, end int
	begin := offsets[0]
	for i, end = range offsets[1:] {
		record[i] = strings.TrimSpace(string(bytes[begin:end]))
		begin = end
	}
	return
}

// NewSepSpliter 创建拆分器
func NewSepSpliter(sep string) SplitFunc {
	return func(s *bufio.Scanner) []string {
		return strings.Split(s.Text(), sep)
	}

}
