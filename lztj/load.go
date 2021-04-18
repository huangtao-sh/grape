package main

import (
	"fmt"
	"grape"
	"grape/loader"
	"grape/sqlite3"
	"strings"
)

func conv(row []string) ([]string, error) {
	s := row[10]
	row[10] = strings.Join([]string{s[:4], s[4:6], s[6:]}, "-")
	return row, nil
}

// LoadZg 导入营业主管信息
func LoadYyzg(file string) {
	info := grape.NewPath(file).FileInfo()
	reader := loader.NewConverter(loader.NewXlsReader(file, 0, 1),
		loader.Include(2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 11),
		conv)
	const loadSQL = "insert or replace into yyzg values(?,?,?,?,?,?,?,?,?,?,?)"
	lder := loader.NewLoader(info, "yyzg", loadSQL, reader)
	lder.Ver = grape.Extract(`\d+`, info.Name())
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
	row[3] = strings.ToUpper(row[3])                  // 将工号全部转换成大写字母
	row[3] = grape.Extract(`[A-Z]{1,2}\d{4}`, row[3]) // 修正员工号，将多余的数字剔除
	if row[3] != "" {
		row = row[:23]
	} else {
		row = nil
	}
	return row, nil
}

// LoadZg 导入营业主管信息
func LoadLzbg(file string) {
	loadSQL := sqlite3.LoadSQL("insert or replace", "lzbg", 23)
	info := grape.NewPath(file).FileInfo()
	reader := loader.NewConverter(loader.NewXlsReader(file, 0, 1), convLzbg)
	lder := loader.NewLoader(info, "lzbg", loadSQL, reader)
	lder.Ver = "1.0"
	lder.Clear = false
	lder.Check = true
	err := lder.Load()
	if err == nil {
		fmt.Printf("%s 导入成功\n", info.Name())
	} else {
		fmt.Println(err)
	}
}
