package main

import (
	"fmt"
	"grape"
	"grape/data/xls"
	"grape/sqlite3"
)

const (
	Version    = "1.0"
	dbsjHeader = "账号,户名,账户性质,开户日期,销户日期,状态"
	dbsjSQL    = `select a.zh,a.hm,a.zhxz,a.khrq,a.xhrq,a.zt from rhsj a
left join bhsj b on a.zh=b.zh 
where b.zh is null and %s 
order by a.zh`
	qssjHeader = "账号,户名,账户类别,开户日期,销户日期,状态"
	qssjSQL    = `select a.zh,a.yshm,a.zhlb,a.khrq,a.xhrq,a.zt from bhsj a left join 
rhsj b on a.zh=b.zh 
where b.zh is null and %s and substr(a.zh,12,3)='201' and a.bz='01' 
order by a.khrq `
	cwsjHeader = "账号,户名,户名（人行）,账户类别,账户性质（人行）,开户日期,开户日期（人行）,销户日期,销户日期（人行）,状态,状态（人行）,户名比对结果,账户种类比对结果,开户日期比对结果,销户日期比对结果,状态比对结果"
	cwsjSQL    = `select a.zh,a.yshm,b.hm,a.zhlb,b.zhxz,a.khrq,b.khrq,a.xhrq,b.xhrq,a.zt,b.zt,
iif(a.hm==b.hm,"","不相符"),iif(a.zhlb==b.zhlb,"","不相符"),iif(a.khrq==b.khrq,"","不相符"),iif(a.xhrq==b.xhrq,"","不相符"),iif(a.zt==b.zt,"","不相符") 
from bhsj a left join rhsj b on a.zh=b.zh  
where(a.hm<>b.hm or a.khrq<>b.khrq or a.xhrq<>b.xhrq or a.zt<>b.zt or a.zhlb<>b.zhlb)and %s
and substr(a.zh,12,3)='201' and a.bz='01' 
order by a.khrq`
	tjHeader = "账号,户名,报送笔数,正确笔数"
	tjSQL    = `select a.zh,a.yshm,count(a.zh),
sum(case when a.hm=b.hm and a.khrq=b.khrq and a.xhrq=b.xhrq and a.zt=b.zt and a.zhlb=b.zhlb then 1 else 0 end)
from bhsj a left join rhsj b on a.zh=b.zh  
where %s
and substr(a.zh,12,3)='201' and a.bz='01' 
group by a.zh having count(a.zh)>1
order by a.zh`
)

// Query test
func Query(exportall bool) {

	var (
		dbsjWidth = map[string]float64{
			"A":   25.4,
			"B":   65,
			"C":   13,
			"D:E": 12,
			"F":   8,
		}
		cwsjWidth = map[string]float64{
			"A":   25.4,
			"B:C": 65,
			"D:E": 13,
			"F:I": 12,
			"J:K": 8,
			"L:P": 15,
		}
		tjWidth = map[string]float64{
			"A":   25.4,
			"B":   65,
			"C:D": 8,
		}
	)
	var condition string
	if exportall {
		condition = "1"
	} else {
		condition = "(a.khrq > date('now','-3 month')or(a.xhrq>date('now','-3 nonth'))) "
	}
	book := xls.NewFile()
	sheet := book.GetSheet(0)
	sheet.Rename("多报送数据")
	sheet.Write("A1", dbsjHeader, dbsjWidth, sqlite3.Fetch(fmt.Sprintf(dbsjSQL, condition)))
	sheet = book.GetSheet("漏报送数据")
	sheet.Write("A1", qssjHeader, dbsjWidth, sqlite3.Fetch(fmt.Sprintf(qssjSQL, condition)))
	sheet = book.GetSheet("报送错误数据")
	sheet.Write("A1", cwsjHeader, cwsjWidth, sqlite3.Fetch(fmt.Sprintf(cwsjSQL, condition)))
	sheet = book.GetSheet("多次报送统计")
	sheet.Write("A1", tjHeader, tjWidth, sqlite3.Fetch(fmt.Sprintf(tjSQL, condition)))
	file := grape.NewPath("~/Downloads/账户报送数据比对.xlsx")
	book.SaveAs(file)
	fmt.Printf("导出报送成功\n")
}
