/*
	项目：内控违规统计程序
	作者：黄涛
	创建：2021-01-09

*/

package main

import (
	"flag"
	"fmt"
	"grape"
	"grape/loader"
	"grape/sqlite3"
	"log"
	"strings"
)

const (
	nksql  = "insert or ignore into nkwg values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	jlrsql = "insert or ignore into djr values(?,?)"
)

func conv(row []string) ([]string, error) {
	return row[:29], nil
}

// Load 读取数据
func Load() {
	file := grape.NewPath("~/Downloads").Find("resultReg*.xls")
	if file == "" {
		fmt.Println("未在 ~/Downloads 目录下找到 resultReg.xls 文件")
	} else {
		p := grape.NewPath(file)
		info := p.FileInfo()
		d := loader.NewXlsReader(file, 0, 1)
		dd := loader.NewConverter(d, conv)
		l := loader.NewLoader(info, "nkwg", nksql, dd)
		l.Clear = false // 不对已导入数据进行清理
		err := l.Load()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("导入文件 %s 成功\n", info.Name())
		}
	}
	file = grape.NewPath("~/Documents/参数备份/考核记录人").Find("考核记录人*.xls")
	if file == "" {
		fmt.Println("未在 ~/Documents/参数备份/考核记录人 目录下找到 考核记录人.xls 文件")
	} else {
		info := grape.NewPath(file).FileInfo()
		d := loader.NewXlsReader(file, 0, 1)
		l := loader.NewLoader(info, "djr", jlrsql, d)
		err := l.Load()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("导入文件 %s 成功\n", info.Name())
		}
	}
}

// getCurMonth 获取最近的月份
func getCurMonth() (year, month string) {
	err := sqlite3.QueryRow("select max(substr(lrsj,1,7)) from nkwg").Scan(&month)
	grape.CheckFatal(err)
	log.Printf("Current month :%s", month)
	year, month = month[:4], month[5:]
	return
}

// Report 报告登记情况
func Report() {
	year, month := getCurMonth()
	value := 1
	if month >= "05" {
		value = 3
	}
	fmt.Printf("当前年份：%s\n", year)
	sqlite3.Printf(
		"合计笔数：%d，合计扣分：%d\n",
		"select count(djbh),sum(kfz) from nkwg where strftime('%Y',lrsj)=?", year)
	fmt.Printf("\n笔数不足%d笔人员清单\n", value)
	fmt.Println("工号    姓名      数量  最后登记日期")
	sqlite3.Printf(
		"%5s  %-10s %4d   %10s\n",
		`select a.gh,a.djr,ifnull(b.sl,0),ifnull(b.sj,"")
		from djr a left join 
		(select lrrgh,count(djbh) as sl,max(lrsj) as sj from nkwg where strftime('%Y',lrsj)=? group by lrrgh)  b
		on a.gh=b.lrrgh  
		where ifnull(sl,0)<?
		order by sl desc
		`, year, value)
}

// ShowAll 显示本部门当前年度扣分情况
func ShowAll() {
	const sql = `
select lrrgh,lxrxm,count(djbh)as sl,sum(kfz),max(lrsj)from nkwg
where strftime('%Y',lrsj)=?
group by lrrgh
order by sl desc`
	const format = "%5s  %-10s  %4d  %4d  %10s\n"
	year, _ := getCurMonth()
	fmt.Printf("当前年份：%s\n", year)
	fmt.Println("工号   姓名        笔数   分值  最后时间")
	sqlite3.Printf(format, sql, year)
	sqlite3.Printf(
		"合计笔数：%d，合计扣分：%d\n",
		"select count(djbh),sum(kfz) from nkwg where strftime('%Y',lrsj)=?", year)
}

func main() {
	grape.InitLog()
	sqlite3.Config("nkwg")
	log.Printf("Set Database %s", "nkwg")
	sqlite3.ExecScript(initSQL)
	load := flag.Bool("l", false, "导入数据")
	report := flag.Bool("r", false, "报告统计结果")
	showall := flag.Bool("a", false, "显示本年度扣分情况")
	query := flag.String("q", "", "查询内控违规")
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
	if *query != "" {
		sql := `
select lrsj,ssjg,sstx,kfz,fxqk from nkwg 
where %s
and lrsj>date('now','-3 year')    -- 仅显示近三年记录
order by lrsj 
`
		var q []string
		for _, v := range strings.Split(*query, " ") {
			q = append(q, fmt.Sprintf("fxqk like '%%%s%%'", v))
		}
		s := strings.Join(q, " and ")
		sqlite3.Printf("录入时间：%s\n所属机构：%s\n所属条线：%s\n扣分值  ：%d\n%s\n\n",
			fmt.Sprintf(sql, s),
			fmt.Sprintf("%%%s%%", *query))
	}
}
