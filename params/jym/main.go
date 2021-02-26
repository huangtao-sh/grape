package jym

import (
	"flag"
	"fmt"
	"grape/data/xls"
	"grape/params"
	"grape/path"
	"grape/sqlite3"
	"grape/util"
)

const (
	header  = `交易名称,交易码,交易组,交易组名称,优先级,网点授权级别,中心授权级别,必须网点授权,中心授权机构,必须中心授权,技能级别,现转标志,是否外包,大额提示,是否扫描电子底卡,是否收手续费,是否需要后台监测,事中扫描方式,补扫时限,是否需要审查,是否允许抹账,是否允许超额授权,附加交易组,事后补扫,磁道校验,一级菜单,二级菜单`
	header2 = `交易名称,交易码,交易组,交易组名称,优先级,网点授权级别,中心授权级别,必须网点授权,中心授权机构,必须中心授权,技能级别,现转标志,是否外包,大额提示,是否扫描电子底卡,是否收手续费,是否需要后台监测,事中扫描方式,补扫时限,是否需要审查,是否允许抹账,是否允许超额授权,附加交易组,事后补扫,磁道校验`
)

// Main jy 程序入口
func Main() {
	defer util.Recover()
	defer sqlite3.Close()
	var (
		export  = flag.Bool("e", false, "导出交易参数表")
		update  = flag.Bool("u", false, "更新参数")
		publish = flag.Bool("p", false, "发布交易参数")
		load    = flag.Bool("l", false, "恢复交易参数")
		backup  = flag.Bool("b", false, "备份交易参数")
		check   = flag.Bool("c", false, "检查交易码是否已使用")
	)
	flag.Parse()
	if *export {
		params.PrintVer("jym")
		exportToXlsx()
	}
	if *publish {
		UpdateJycs() // 更新交易码参数
		Publish()    // 发布交易码参数
	}
	if *load {
		LoadJycs() // 导入交易参数
	}
	if *backup {
		BackupJycs() // 备份交易参数
	}
	if *update {
		UpdateJycs() // 更新交易码参数
	}
	if len(flag.Args()) > 0 {
		const Format = "%-50s  %4s  %6s  %-20s\n"
		params.PrintVer("jym")
		for _, arg := range flag.Args() {
			if util.FullMatch(`\d{4}`, arg) {
				if *check {
					fmt.Println("检查交易码占用情况")
					sqlite3.Printf("生产参数：%4s  %-40s\n", "select jym,jymc from jym where jym=?", arg)
					sqlite3.Printf("待投产  ：%4s  %-40s\n", "select jym,jymc from jymcs where jym=?", arg)
					sqlite3.Printf("交易菜单：%4s  %-40s\n", "select jym,name from menu where jym=?", arg)
				} else {
					err := sqlite3.PrintRow(header, "select * from jycs where jym=?", arg)
					if err != nil {
						fmt.Printf("错误：交易码 %s 不存在\n", arg)
					} else {
						var tcrq string
						sqlite3.QueryRow("select tcrq from jymcs where jymc=?", arg).Scan(&tcrq)
						if tcrq != "" {
							fmt.Printf("投产日期：%s", tcrq)
						}
					}
				}
			} else if util.FullMatch(`[A-Z]{2}\d{3}[A-Z]{1}`, arg) {
				sqlite3.Printf(Format, "select jymc,jym,jyz,jyzm from jycs where jyz=? order by jym", arg)
			} else {
				sqlite3.Printf(Format, "select jymc,jym,jyz,jyzm from jycs where jymc like ? order by jym", fmt.Sprintf(`%%%s%%`, arg))
			}
		}
	}
}

func exportToXlsx() {
	const (
		queryJycs = `select * from jycs order by jym`
		queryJymb = `select * from jymb order by jym`
		jySheet   = "交易码表"
		csSheet   = "交易码参数表"
	)
	var (
		jymbWidth = map[string]float64{
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
		jycsWidth = map[string]float64{
			"A":   44,
			"B:C": 7,
			"D":   21,
			"E":   7,
			"F:V": 13,
			"W":   22,
			"X:Y": 9,
		}
	)
	book := xls.NewFile()
	sheet := book.GetSheet(0)
	sheet.Rename(jySheet)
	sheet.Write("A1", header, jymbWidth, sqlite3.Fetch(queryJycs))
	sheet = book.GetSheet(csSheet)
	sheet.Write("A1", header2, jycsWidth, sqlite3.Fetch(queryJymb))
	file := path.Home.Join("Documents", fmt.Sprintf("交易码参数表（%s）.xlsx", params.GetVer("jym")))
	book.SaveAs(file)
	fmt.Printf("导出参数成功")
}
