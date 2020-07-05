package jym

import (
	"flag"
	"fmt"
	"grape/data/xls"
	"grape/params"
	"grape/path"
	"grape/sqlite3"
	"grape/text"
	"regexp"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

var header = `交易名称,交易码,交易组,交易组名称,优先级,网点授权级别,中心授权级别,必须网点授权,中心授权机构,必须中心授权,技能级别,现转标志,是否外包,大额提示,是否扫描电子底卡,是否收手续费,是否需要后台监测,事中扫描方式,补扫时限,是否需要审查,是否允许抹账,是否允许超额授权,附加交易组,事后补扫,磁道校验,一级菜单,二级菜单`

// Main jy 程序入口
func Main() {
	params.PrintVer("jym")
	Fmt := "%-50s  %4s  %6s  %-20s\n"
	var export = flag.Bool("e", false, "导出交易参数表")
	flag.Parse()
	if *export {
		exportToXlsx()
	}
	for _, arg := range flag.Args() {
		if matched, _ := regexp.MatchString(`^\d{4}$`, arg); matched {
			err := sqlite3.PrintRow(header, "select * from jycs where jym=?", arg)
			if err != nil {
				fmt.Printf("错误：交易码 %s 不存在\n", arg)
			}
		} else if matched, _ := regexp.MatchString(`^[A-Z]{2}\d{3}[A-Z]{1}$`, arg); matched {
			sqlite3.Printf(Fmt, "select jymc,jym,jyz,jyzm from jycs where jyz=?", arg)
		} else {
			sqlite3.Printf(Fmt, "select jymc,jym,jyz,jyzm from jycs where jymc like ?", fmt.Sprintf(`%%%s%%`, arg))
		}
	}
}

var jymbWidth = map[string]float64{
	"A":   44,
	"B:C": 7,
	"D":   21,
	"E":   7,
	"F:V": 13,
	"W":   22,
	"X:Y": 9,
	"Z":   17,
	"AA":  33,
}
var jycsWidth = map[string]float64{
	"A":   44,
	"B:C": 7,
	"D":   21,
	"E":   7,
	"F:V": 13,
	"W":   22,
	"X:Y": 9,
}

func exportToXlsx() {
	filename := fmt.Sprintf("交易码参数表（%s）.xlsx", params.GetVer("jym"))
	filename = (path.NewPath("~/Documents").Join(filename)).String()
	book := excelize.NewFile()
	book.NewSheet("交易码表")
	xls.SetWidth(book, "交易码表", jymbWidth)
	xls.SetWidth(book, "交易码参数表", jycsWidth)
	book.NewSheet("交易码参数表")
	book.DeleteSheet("Sheet1")
	head := text.Slice(strings.Split(header, ","))
	book.SetSheetRow("交易码表", "A1", &head)
	rows := sqlite3.Fetch("select * from jycs order by jym")
	rows.Export(book, "交易码表", "A2")
	head = head[:len(head)-2]
	book.SetSheetRow("交易码参数表", "A1", &head)
	rows = sqlite3.Fetch("select * from jymb order by jym")
	rows.Export(book, "交易码参数表", "A2")
	book.SaveAs(filename)
	fmt.Printf("导出参数成功")
}
