package gbk

import (
	"testing"
)

func TestEncode(t *testing.T) {
	s := "中文编码测试"
	b, _ := Encode(s)
	k, _ := Decode(b)
	if s != k {
		t.Errorf("测试失败")
	}
}
