package gbk

import (
	"bytes"
	"io/ioutil"
	"testing"
)

var (
	s = "中文编码测试 abc"
	b = []byte("\xd6\xd0\xce\xc4\xb1\xe0\xc2\xeb\xb2\xe2\xca\xd4 abc")
)

func TestReader(t *testing.T) {
	r := NewReader(bytes.NewReader(b))
	s2, err := ioutil.ReadAll(r)
	if err != nil {
		t.Errorf("Test Failed")
	}
	if s != string(s2) {
		t.Errorf("Test Failed")
	}
}

func TestWriter(t *testing.T) {
	bs := make([]byte, 0)
	w1:=bytes.NewBuffer(bs)
	w := NewWriter(w1)
	w.Write([]byte(s))
	b2:=w1.String()
	if b2!=string(b){
		t.Errorf("Test Failed")
	}
}

func TestEncode(t *testing.T) {
	s := "中文编码测试 abc"
	b := []byte("\xd6\xd0\xce\xc4\xb1\xe0\xc2\xeb\xb2\xe2\xca\xd4 abc")
	b1, err := Encode(s)
	if err != nil {
		t.Errorf("编码测试失败")
	}
	if string(b) != string(b1) {
		t.Errorf("编码测试失败")
	}
	s1, err := Decode(b)
	if err != nil {
		t.Errorf("解码测试失败")
	}
	if s != s1 {
		t.Errorf("解码测试失败")
	}
}
