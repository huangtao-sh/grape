package wwxt

import (
	"fmt"
	"grape/data/xls"
	"grape/params/load"
	"grape"
	"io"
	"os"
	"strconv"
)

const initWwxt = `
create table if not exists wwxt(
	id 	int primary key,
	name	text,
	jglx	text  default '3-总分全辖',
	jgh		text  default '999999999',
	date	text  
)
`
const loadWwxt = `insert into wwxt values(?,?,?,?,?)`

func conv(row []string) []string {
	_, err := strconv.Atoi(row[0])
	if len(row) < 5 {
		row = append(row, "")
	}
	if err == nil {
		row[4] = xls.ConvertDate(row[4])
	} else {
		row = nil
	}
	return row
}

// LoadWwxt 导入外围系统
func LoadWwxt(info os.FileInfo, r io.Reader, ver string) *load.Loader {
	reader := xls.NewXlsReader(r, "历史", 1, conv)
	return load.NewLoader("wwxt", info, ver, reader, initWwxt, loadWwxt)
}

// Load 导入数据主程序
func Load() {
	file := Root.Find("新增外围系统列表????-??-??.xlsx")
	if file != "" {
		p := grape.NewPath(file)
		r, err := p.Open()
		ver := p.FileInfo().Name()[24:34]
		grape.CheckFatal(err)
		loader := LoadWwxt(p.FileInfo(), r, ver)
		loader.Load()
	} else {
		fmt.Print("未发现文件")
	}
}
