package xyk

import (
	"strings"
	"strconv"
)

func Atoi(s string)(result int){
	s=strings.Replace(s,".","",1)
	s=strings.Replace(s,"+","",1)
	result,err:=strconv.Atoi(s)
	if err!=nil{
		panic("转换数字失败")
	}
	return
}