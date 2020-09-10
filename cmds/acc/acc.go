package main

import (
	"flag"
	"fmt"
	"grape/params"
	"grape/sqlite3"
	"strings"
)

// 打印该科目最小未用序号
func minAvaible(kemu string) {
	var i, k int
	r := sqlite3.Fetch("select distinct xh from nbzhhz where km=? order by xh", kemu)
	defer r.Close()
	for i = 1; r.Next(); i++ {
		r.Scan(&k)
		if i != k {
			break
		}
	}
	fmt.Printf("\n最小未用序号： %03d\n\n", i)
}

// 打印该科目开户情况
func showKaihu(kemu string) {
	fmt.Print("\n已开账户情况")
	sqlite3.Printf("%s-%03d %s\n", "select distinct km,xh,hm from nbzhhz where km=? order by xh", kemu)
}

// 打印该科目参数
func kmcs(kemu string) {
	const (
		queryKemu   = `select * from kemu where km=?`
		queryFormat = "%s  %s\n%s\n"
	)
	sqlite3.Printf(queryFormat, queryKemu, kemu)
	fmt.Println("\n科目参数")
	sqlite3.PrintRow(
		"科    目,汇总科目,科目名称,科目级别,借贷标志,科目类型",
		`select km,hzkm,kmmc,
		case kmjb when '0' then'0-明细科目' when '1' then '1-一级科目' when '2' then '2-二级科目' when '3' then '3-三级科目' end, 
		case jdbz when '0' then '0-两性' when '1' then '1-借方' when '2' then '2-贷方' when '3' then '3-并列 (借贷不轧差)' end,
		case kmlx when '0' then '0-汇总科目 (不开户)' when '1' then '1-单账户科目' when '2' then '2-多账户科目' end
		from kmzd where km=?
		`, kemu)
}

// 科目查询主程序
func main() {
	defer sqlite3.Close()  // 程序结束时，关闭数据库
	var km, xh string
	flag.Parse()
	params.PrintVer("nbzhmb") // 打印参数版本
	for _, ac := range flag.Args() {
		kk := strings.Split(ac, "-")
		if len(kk) == 2 {
			km, xh = kk[0], kk[1]
		} else if len(ac) == 6 {
			kmcs(ac)
			showKaihu(ac)
			minAvaible(ac)
			fmt.Println("内部账户开立模板")
			fmt.Println("维护日期     机构   科目  序号  币种           透支额度   状态  计息  户名")
			sqlite3.Printf("%-10s  %4s  %6s  %03d  %4s  %19,.2f  %4s  %4s  %s\n",
				"select whrq,jglx,km,xh,bz,tzed,cszt,jxbz,hm from nbzhmb where km=? order by km,xh,jglx,bz", ac)
			return
		} else if len(ac) == 9 {
			km, xh = ac[:6], ac[6:]
		} else {
			fmt.Printf("账号:%s 格式错，应为：000000-1或者：000000000\n", ac)
			return
		}
		fmt.Printf("科目：%s  序号：%s\n", km, xh)
		fmt.Println("维护日期     机构   科目  序号  币种           透支额度   状态  计息  户名")
		sqlite3.Printf("%-10s  %4s  %6s  %03d  %4s  %19,.2f  %4s  %4s  %s\n",
			"select whrq,jglx,km,xh,bz,tzed,cszt,jxbz,hm from nbzhmb where km=? and xh=?", km, xh)
	}
}
