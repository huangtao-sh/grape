package wwxt

import (
	"fmt"
	"grape/data/xls"
	"grape/params/load"
	"grape/path"
	"grape/util"
	"io"
	"os"
	"strconv"
	"strings"
)

var initWwxt = `
create table if not exists wwxt(
	id 	int primary key,
	name	text,
	jglx	text  default '3-总分全辖',
	jgh		text  default '999999999',
	date	text  
)
`
var loadWwxt = `insert into wwxt values(?,?,?,?,?)`

func conv(row []string) []string {
	_, err := strconv.Atoi(row[0])
	if len(row) < 5 {
		row = append(row, "")
	}
	if err == nil {
		row[4] = ConvDate(row[4])
	} else {
		row = nil
	}
	return row
}

// LoadWwxt 导入外围系统
func LoadWwxt(info os.FileInfo, r io.Reader, ver string) *load.Loader {
	reader := xls.NewXlsReader(r, "历史", 1, conv)
	return load.NewLoader("jyz", info, ver, reader, initWwxt, loadWwxt)
}

// ConvDate 转换日期，把 05-16-20 格式的日期转换成 2020-05-16 格式，无法转换则返回原数据
func ConvDate(d string) string {
	if len(d) == 8 {
		s := strings.Split(d, "-")
		return fmt.Sprintf("20%s-%s-%s", s[2], s[0], s[1])
	}
	return d
}

// Load 导入数据主程序
func Load() {
	file := Root.Find("新增外围系统列表????-??-??.xlsx")
	if file != "" {
		p := path.NewPath(file)
		r, err := p.Open()
		ver := p.FileInfo().Name()[24:34]
		util.CheckFatal(err)
		loader := LoadWwxt(p.FileInfo(), r, ver)
		loader.Load()
	} else {
		fmt.Print("未发现文件")
	}
}
