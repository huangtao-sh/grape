package text

import "strings"

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

// SepSpliter 按指定分割符拆分
type SepSpliter struct {
	sep string
}

// NewSepSpliter 根据分割符拆分
func NewSepSpliter(sep string) *SepSpliter {
	return &SepSpliter{sep}
}

// Split 拆分
func (s *SepSpliter) Split(bytes []byte) []interface{} {
	return Slice(strings.Split(string(bytes), s.sep))
}
