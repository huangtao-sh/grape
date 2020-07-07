package jym

import (
	"grape/data/xls"
	"grape/params/load"
	"io"
	"os"
)

const loadJyz = `insert into jyz values(?,?)`

func conv(s []string) []string {
	if s[0] != "" {
		return s
	}
	return nil
}

// LoadJyz 导入交易组
func LoadJyz(info os.FileInfo, r io.Reader, ver string) *load.Loader {
	InitJym()
	reader := xls.NewXlsReader(r, "交易组", 1, conv)
	return load.NewLoader("jyz", info, ver, reader, "", loadJyz)
}
