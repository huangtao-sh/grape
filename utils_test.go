package grape

import (
	"testing"
)

func TestExtractPos(t *testing.T) {
	if ExtractPos(`(\d{3})`, "fsdab134fsda", 1) != "134" {
		t.Error("test ExtractPos failed")
	}
	if ExtractPos(`(.*?银行|.*?公司)?(.*?分行)`, "上海银行重庆分行", 2) != "重庆分行" {
		t.Error("test ExtractPos failed")
	}
	if ExtractPos(`(.*?公司|.*?银行)?(.*?分行)`, "上海银行股份有限公司重庆分行", 2) != "重庆分行" {
		t.Error("test ExtractPos failed")
	}
}

func TestSprintf(t *testing.T){
	s1:="重庆"
	s2:="重庆  "
	if Sprintf("%-6s",s1)!=s2{
		t.Error("test Sprintf failed")
	}
	if Sprintf("%6s",s1)!="  "+s1{
		t.Error("test Sprintf failed")
	}
}

func TestWlen(t *testing.T){
	if Wlen("河南省123")!=9{
		t.Error("test Wlen failed")
	}
}