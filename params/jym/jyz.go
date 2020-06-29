package jym

import (
	"grape/data/xls"
	"grape/params/load"
	"io"
	"os"
)

var initJyz = `
create table if not exists jyz(
	jyz 	text 	priamry key, -- 交易组
	jyzm	text	-- 交易组名
)
`

var loadJyz = `insert into jyz values(?,?)`

func conv(s []string) []string {
	if s[0] != "" {
		return s
	} else {
		return nil
	}
}

// LoadJyz 导入交易组
func LoadJyz(info os.FileInfo, r io.Reader, ver string) *load.Loader {
	reader := xls.NewXlsReader(r, "交易组", 1, conv)
	return load.NewLoader("jyz", info, ver, reader, initJyz, loadJyz)
}
