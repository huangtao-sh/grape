package jym

import (
	"grape/params/load"
	"grape/text"
	"io"
	"os"
)

var loadJymSQL = `insert or replace into jym values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

func convert(s []string) []string {
	if len(s) < 22 {
		s = append(s, "")
	}
	return s
}

// LoadJym 导入交易码文件
func LoadJym(info os.FileInfo, r io.Reader, ver string) *load.Loader {
	InitJym()
	reader := text.NewReader(r, false, text.NewSepSpliter(","), convert)
	return load.NewLoader("jym", info, ver, reader, "", loadJymSQL)
}

var loadShbs = `insert or replace into shbs values(?)`

// LoadShbs 导入交易码文件
func LoadShbs(info os.FileInfo, r io.Reader, ver string) *load.Loader {
	InitJym()
	reader := text.NewReader(r, false, text.NewSepSpliter(","), text.Include(0))
	return load.NewLoader("shbs", info, ver, reader, "", loadShbs)
}

var loadCdjy = `insert or replace into cdjy values(?)`

func convCdjy(s []string) (d []string) {
	if s[1] == "8" {
		d = append(d, s[0])
	} else {
		d = nil
	}
	return
}

// LoadCdjy 导入交易码文件
func LoadCdjy(info os.FileInfo, r io.Reader, ver string) *load.Loader {
	InitJym()
	reader := text.NewReader(r, false, text.NewSepSpliter(","), convCdjy)
	return load.NewLoader("cdjy", info, ver, reader, "", loadCdjy)
}
