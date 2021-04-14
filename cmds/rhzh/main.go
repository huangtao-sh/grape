package main

import (
	"flag"
	"fmt"
	"grape/path"
	"grape/rhzh"
	"grape/sqlite3"
	"os"
)

const (
	description = `人行报送数据比对程序 Version %s
用法：
1、将从人行账户管理系统导出的文件 "单位银行结算账户开立、变更及撤销情况查询结果输出" 放在当前用户的下载目录下
2、将从柜面系统 1181-开销户登记簿中导出的数据文件 "开销户登记簿" 放在当前用户的下载目录下
3、运行本程序后将在 下载 目录生成 "账户报送数据比对" 文件
4、默认仅比对近三个月数据，如需要全量比对，请运行：rhzh -a`
)

func main() {
	path.InitLog()
	version := flag.Bool("v", false, "显示程序版本")
	all := flag.Bool("a", false, "比对所有账户")
	querySQL := flag.String("q", "", "执行查询")
	flag.Parse()
	if *querySQL != "" {
		sqlite3.Println(*querySQL)
		os.Exit(0)
	}
	if *version {
		fmt.Printf(description, rhzh.Version)
		os.Exit(0)
	}
	if 1 == 2 {

		rhzh.LoadBhsj()
		rhzh.LoadRhsj()
	}

	rhzh.Query(*all)
}
