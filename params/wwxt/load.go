package main

import (
	"fmt"
	"grape/data"
	"grape/params"
	"grape/path"
	"grape/sqlite3"
	"grape/util"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// Init 初始化数据库
func Init() {
	sqlite3.ExecScript(`
-- drop table if exists wwxt;
create table if not exists wwxt(
	id 	int primary key,
	name	text,
	jglx	text  default '3-总分全辖',
	jgh		text  default '999999999',
	date	text  
)
	`)
}

// ConvDate 转换日期，把 05-16-20 格式的日期转换成 2020-05-16 格式，无法转换则返回原数据
func ConvDate(d string) string {
	if len(d) == 8 {
		s := strings.Split(d, "-")
		return fmt.Sprintf("20%s-%s-%s", s[2], s[0], s[1])
	}
	return d
}

// LoadFile
func LoadFile(filename string) {
	file := path.NewPath(filename)
	err := sqlite3.ExecTx(
		params.NewChecker("wwxt", file, file.Base()),
		NewFile(filename),
	)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(file.Base(), "导入完成")
	}
}

// Load 导入数据主程序
func Load() {
	Init()                                        // 初始化数据库
	files := Root.Glob("新增外围系统列表????-??-??.xlsx") // 查找数据文件
	if len(files) > 0 {
		LoadFile(files[len(files)-1]) // 找到文件，并执行导入操作
	} else {
		fmt.Print("未发现文件")
	}
}

type File struct {
	data.Data
	file string
}

func NewFile(file string) *File {
	return &File{*data.NewData(), file}
}

func (f *File) Read() {
	defer f.Close()
	xls, err := excelize.OpenFile(f.file)
	util.CheckFatal(err)
	rows, err := xls.Rows("历史")
	util.CheckFatal(err)
	ch := f.WriteCh()
	for rows.Next() {
		row, _ := rows.Columns()
		id, err := strconv.Atoi(row[0])
		if len(row) < 5 {
			row = append(row, "")
		}
		if err == nil {
			ch <- []interface{}{id, row[1], row[2], row[3], ConvDate(row[4])}
		}
	}
}
func (f *File) Exec(tx *sqlite3.Tx) error {
	tx.Exec("delete from wwxt") // 清空现有数据
	f.Add(1)
	go f.Data.Exec(tx, "insert or replace into wwxt values(?,?,?,?,?)")
	//go tx.ExecCh("insert or replace into wwxt values(?,?,?,?,?)", &f.Data)
	go f.Read()
	f.Wait()
	return nil
}
