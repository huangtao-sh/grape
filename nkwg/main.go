package nkwg

import (
	"flag"
	"fmt"
	"grape/data"
	"grape/params/load"
	"grape/path"
	"grape/sqlite3"
	"grape/util"
)

// Load 导入数据
func Load() {
	sql := "insert or ignore into nkwg values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	jlrsql := util.Sprintf("insert or ignore into djr %2V")
	file := path.NewPath("~/Downloads").Find("resultReg*.xls")
	if file == "" {
		fmt.Println("未在 ~/Downloads 目录下找到 resultReg.xls 文件")
	} else {
		p := path.NewPath(file)
		info := p.FileInfo()
		r := data.NewXlsReader(file, 0, 1)
		d := data.NewConvertReader(r)
		loader := load.NewLoader("nkwg", info, "", d, "", sql)
		loader.Load()
	}
	file = path.NewPath("~/OneDrive/工作/参数备份/考核记录人").Find("考核记录人*.xls")
	if file == "" {
		fmt.Println("未在 ~/OneDrive/工作/参数备份/考核记录人 目录下找到 考核记录人.xls 文件")
	} else {
		info := path.NewPath(file).FileInfo()
		r := data.NewXlsReader(file, 0, 1)
		d := data.NewConvertReader(r)
		loader := load.NewLoader("djr", info, "", d, "", jlrsql)
		loader.Load()
	}
}

// Report 报告登记情况
func Report() {
	var year string
	sqlite3.QueryRow("select max(substr(lrsj,1,7)) from nkwg").Scan(&year)
	month := string([]byte(year)[5:])
	value := 1
	if month >= "05" {
		value = 3
	}
	year = string([]byte(year)[:4])
	fmt.Printf("当前年份：%s\n", year)
	sqlite3.Printf(
		"合计笔数：%d，合计扣分：%d\n",
		"select count(djbh),sum(kfz) from nkwg where strftime('%Y',lrsj)=?", year)
	fmt.Printf("\n笔数不足%d笔人员清单\n", value)
	fmt.Println("工号    姓名      数量  最后登记日期")
	sqlite3.Printf(
		"%5s  %-10s %4d   %10s\n",
		`select a.gh,a.djr,nullif(b.sl,0) as sl,nullif(b.sj,"")
		from djr a left join 
		(select lrrgh,count(djbh) as sl,max(lrsj) as sj from nkwg where strftime('%Y',lrsj)=? group by lrrgh)  b
		on a.gh=b.lrrgh  
		where sl<?
		order by sl desc
		`, year, value)
}

// ShowAll 显示本年度扣分情况
func ShowAll() {
	const sql = `
select lrrgh,lxrxm,count(djbh)as sl,sum(kfz),max(lrsj)from nkwg
where strftime('%Y',lrsj)=?
group by lrrgh
order by sl desc`
	const format = "%5s  %-10s  %4d  %4d  %10s\n"
	var year string
	sqlite3.QueryRow("select max(strftime('%Y',lrsj)) from nkwg").Scan(&year)
	fmt.Printf("当前年份：%s\n", year)
	fmt.Println("工号   姓名        笔数   分值  最后时间")
	sqlite3.Printf(format, sql, year)
	sqlite3.Printf(
		"合计笔数：%d，合计扣分：%d\n",
		"select count(djbh),sum(kfz) from nkwg where strftime('%Y',lrsj)=?", year)
}

// Main 主程序
func Main() {
	load := flag.Bool("l", false, "导入数据")
	report := flag.Bool("r", false, "报告统计结果")
	showall := flag.Bool("a", false, "显示本年度扣分情况")

	flag.Parse()
	if *load {
		Load()
	}
	if *report {
		Report()
	}
	if *showall {
		ShowAll()
	}
}
