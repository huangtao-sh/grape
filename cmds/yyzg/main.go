package main

import (
	"flag"
	"fmt"
	"grape/sqlite3"
	"grape/util"
)

func main() {
	sqlite3.Config("lzbg")
	defer sqlite3.Close()
	var t string
	typ := flag.String("t", "", "主管类型")
	flag.Parse()
	if *typ != "" {
		t = fmt.Sprintf(`and js like "%s%%" `, *typ)
	}
	var ver string
	err := sqlite3.QueryRow("select ver from loadfile where name='yyzg' ").Scan(&ver)
	util.CheckFatal(err)
	fmt.Printf("数据版本：%s\n", ver)

	fmt.Println("工号   姓名       角色            联系电话           手机         机构")
	for _, arg := range flag.Args() {
		if util.FullMatch(`\d{4,9}`, arg) {
			arg = fmt.Sprintf("%s%%", arg)
			sqlite3.Printf("%-6s %-10s %-15s %-15s %11s %-30s\n",
				"select ygh,xm,js,lxdh,mobile,jgmc from yyzg where jg like ? "+t, arg)
		} else {
			arg = fmt.Sprintf("%%%s%%", arg)
			sqlite3.Printf("%-6s %-10s %-15s %-15s %11s %-30s\n",
				"select ygh,xm,js,lxdh,mobile,jgmc from yyzg where(xm like ? or jgmc like ?)"+t, arg, arg)
		}
	}
}
