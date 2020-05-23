package main

import (
	"flag"
	"fmt"
	"grape/params"
	"grape/sqlite3"
	"strings"
)

func main() {
	var km, xh string
	flag.Parse()
	params.PrintVer("nbzhmb")
	for _, ac := range flag.Args() {
		kk := strings.Split(ac, "-")
		if len(kk) == 2 {
			km, xh = kk[0], kk[1]
		} else if len(ac) == 6 {
			fmt.Printf("科目：%s\n", ac)
			fmt.Println("维护日期     机构   科目  序号  币种         透支额度   状态  计息  户名")
			sqlite3.Printf("%-10s  %4s  %6s  %03d  %4s  %17.2f  %4s  %4s  %s\n", 
				"select whrq,jglx,km,xh,bz,tzed,cszt,jxbz,hm from nbzhmb where km=? order by km,xh,jglx,bz", ac)
			return
		} else if len(ac) == 9 {
			km, xh = ac[:6], ac[6:]
		} else {
			fmt.Printf("账号:%s 格式错，应为：000000-1或者：000000000\n", ac)
			return
		}
		fmt.Printf("科目：%s  序号：%s\n", km, xh)
		fmt.Println("维护日期     机构   科目  序号  币种         透支额度   状态  计息  户名")
		sqlite3.Printf("%-10s  %4s  %6s  %03d  %4s  %17.2f  %4s  %4s  %s\n",
			"select whrq,jglx,km,xh,bz,tzed,cszt,jxbz,hm from nbzhmb where km=? and xh=?", km, xh)
	}
}
