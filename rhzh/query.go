package rhzh

import (
	"fmt"
	"grape/data/xls"
	"grape/path"
	"grape/sqlite3"
)

const (
	dbsjHeader = "账号,户名,账户性质,开户日期,销户日期,状态"
	dbsjSQL    = `select a.zh,a.hm,a.zhxz,a.khrq,a.xhrq,a.zt from rhsj a
left join bhsj b on a.zh=b.zh or a.zh="NRA"||b.zh
where b.zh is null
order by a.zh`
	qssjHeader = "账号,户名,账户类别,开户日期,销户日期,状态"
	qssjSQL    = `select a.zh,a.yshm,a.zhlb,a.khrq,a.xhrq,a.zt from bhsj a left join 
rhsj b on a.zh=b.zh or "NRA"||a.zh=b.zh 
where b.zh is null 
order by a.khrq `
	cwsjHeader = "账号,户名,户名（人行）,账户类别,账户性质（人行）,开户日期,开户日期（人行）,销户日期,销户日期（人行）,状态,状态（人行）"
	cwsjSQL    = `select a.zh,a.yshm,b.hm,a.zhlb,b.zhxz,a.khrq,b.khrq,a.xhrq,b.xhrq,a.zt,b.zt 
from bhsj a left join rhsj b on a.zh=b.zh or "NRA"||a.zh=b.zh 
where a.hm<>b.hm or a.khrq<>b.khrq or a.xhrq<>b.xhrq or a.zt<>b.zt or a.zhlb<>b.zhlb 
order by a.khrq`
	tjHeader = "账号,户名,报送笔数,正确笔数"
	tjSQL    = `select a.zh,a.yshm,count(a.zh),
sum(case when a.hm=b.hm and a.khrq=b.khrq and a.xhrq=b.xhrq and a.zt=b.zt and a.zhlb=b.zhlb then 1 else 0 end)
from bhsj a left join rhsj b on a.zh=b.zh or "NRA"||a.zh=b.zh 
group by a.zh having count(a.zh)>1
order by a.zh`
)

// Query test
func Query() {

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
		}
		tjWidth = map[string]float64{
			"A":   25.4,
			"B":   65,
			"C:D": 8,
		}
	)

	book := xls.NewFile()
	sheet := book.GetSheet(0)
	sheet.Rename("多报送数据")
	sheet.Write("A1", dbsjHeader, dbsjWidth, sqlite3.Fetch(dbsjSQL))
	sheet = book.GetSheet("漏报送数据")
	sheet.Write("A1", qssjHeader, dbsjWidth, sqlite3.Fetch(qssjSQL))
	sheet = book.GetSheet("报送错误数据")
	sheet.Write("A1", cwsjHeader, cwsjWidth, sqlite3.Fetch(cwsjSQL))
	sheet = book.GetSheet("多次报送统计")
	sheet.Write("A1", tjHeader, tjWidth, sqlite3.Fetch(tjSQL))

	file := path.NewPath("~/Downloads/账户报送数据比对.xlsx")
	book.SaveAs(file)
	fmt.Printf("导出报送成功\n")

}
