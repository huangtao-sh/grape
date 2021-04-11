package lzbg

import (
	"fmt"
	"grape/loader"
	"grape/path"
	"grape/util"
	"strings"
)

const loadSQL = "insert or replace into yyzg values(?,?,?,?,?,?,?,?,?,?,?)"

func conv(row []string) ([]string, error) {
	s := row[10]
	row[10] = strings.Join([]string{s[:4], s[4:6], s[6:]}, "-")
	return row, nil
}

// LoadZg 导入营业主管信息
func LoadYyzg(file string) {
	info := path.NewPath(file).FileInfo()
	reader := loader.NewConverter(loader.NewXlsReader(file, 0, 1),
		loader.Include(2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 11),
		conv)
	lder := loader.NewLoader(info, "yyzg", loadSQL, reader)
	lder.Ver = util.Extract(`\d+`, info.Name())
	lder.Clear = true
	lder.Check = true
	err := lder.Load()
	if err == nil {
		fmt.Printf("%s 导入成功\n", info.Name())
	} else {
		fmt.Println(err)
	}
}

func convLzbg(row []string) ([]string, error) {
	row[3] = strings.ToUpper(row[3])
	row[3] = util.Extract(`[A-Z]{1,2}\d{4}`, row[3])
	if row[3] != "" {
		row = row[:23]
	} else {
		row = nil
	}
	return row, nil
}

// LoadZg 导入营业主管信息
func LoadLzbg(file string) {
	loadSQL := util.Sprintf("insert or replace into lzbg %23V")
	info := path.NewPath(file).FileInfo()
	reader := loader.NewConverter(loader.NewXlsReader(file, 0, 1), convLzbg)
	lder := loader.NewLoader(info, "lzbg", loadSQL, reader)
	lder.Ver = "1.0"
	lder.Check = true
	err := lder.Load()
	if err == nil {
		fmt.Printf("%s 导入成功\n", info.Name())
	} else {
		fmt.Println(err)
	}
}
