package main

import (
	"fmt"
	"grape/data/xls"
	"grape/sqlite3"
	"log"
)

const (
	tjSQL = `select bglx,count(distinct ygh) 
from lzbg 
where bgrq>=? and bglx in ('营业主管','事后监督')
group by bglx`
	shjdSQL = `
select a.jg,a.jgmc
from 
(select distinct jg,jgmc from yyzg where jg like '%000' 
and jg not in('331000000','338702000') ) a
left join (select distinct bgr,jg from lzbg where bgrq>=? and bglx='事后监督') b
on instr(a.jgmc,b.jg)
where b.bgr is null 
order by a.jg
`
	shjdHeader = `机构号,机构名称`
	yyzgLbSQL  = `
select a.jg,a.jgmc,a.ygh,a.xm
from 
(select jg,jgmc,ygh,xm from yyzg 
where jg not like "331000%" and js like "a%" 
and jg not in("191000000","342002000","361000000","421000000","551000000","338702000")
) a
left join (select distinct ygh from lzbg where bgrq>=? and bglx='营业主管') b
on a.ygh=b.ygh
where b.ygh is null
order by a.jg
`
	yyzgLbHeader = `机构号,机构名称,员工号,姓名`
	yyzgYcSQL    = `
select a.jg||a.bm,a.ygh,a.bgr from 
(select distinct ygh,bgr,jg,bm from lzbg where bgrq>=? and bglx='营业主管') a
left join
(select ygh from yyzg where js like "a%" and jg not like "331000%" )b 
on a.ygh=b.ygh
where b.ygh is null
order by a.jg||a.bm
`
	yyzgYcHeader = `机构,员工号,报告人`
)

func report() {
	var current string
	err := sqlite3.QueryRow(`select date(max(bgrq),"start of month","-5 day")from lzbg`).Scan(&current)
	if err != nil {
		fmt.Println("履职报告中无数据，请先导入数据")
	}
	fmt.Printf("\n起始日期：%s\n", current)
	sqlite3.Printf("%-10s        %3d\n", tjSQL, current)
	fmt.Println("\n事后监督报告漏报清单")
	sqlite3.Printf("%-9s       %-40s\n", shjdSQL, current)
	fmt.Println("\n营业主管履职报告漏报清单")
	sqlite3.Printf("%9s    %-40s      %-6s     %-10s\n", yyzgLbSQL, current)
	fmt.Println("\n营业主管履职报告异常清单")
	sqlite3.Printf("%-40s    %-6s      %-20s\n", yyzgYcSQL, current)

	book := xls.NewFile()
	shjdWidth := map[string]float64{
		"A": 12,
		"B": 50,
	}
	sheet := book.GetSheet(0)
	sheet.Rename("事后监督漏报清单")
	sheet.Write("A1", shjdHeader, shjdWidth, sqlite3.Fetch(shjdSQL, current))

	yyzgLbWidth := map[string]float64{
		"A":   12,
		"B":   50,
		"C:D": 15,
	}

	sheet = book.GetSheet("营业主管漏报清单")
	sheet.Write("A1", yyzgLbHeader, yyzgLbWidth, sqlite3.Fetch(yyzgLbSQL, current))

	yyzgYcWidth := map[string]float64{
		"A":   50,
		"B:C": 15,
	}

	sheet = book.GetSheet("营业主管异常报送清单")
	sheet.Write("A1", yyzgYcHeader, yyzgYcWidth, sqlite3.Fetch(yyzgYcSQL, current))

	book.SaveAs("~/Downloads/履职报告报送统计.xlsx")
	log.Println("导出文件：~/Downloads/履职报告报送统计.xlsx")
}
