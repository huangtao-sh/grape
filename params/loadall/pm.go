package loadall

import (
	"flag"
	"fmt"
	"grape/sqlite3"
)

// Main pm 程序主函数
func Main() {
	defer sqlite3.Close()
	execSQL := flag.String("e", "", "执行 SQL 语句")
	querySQL := flag.String("q", "", "执行查询")
	load := flag.Bool("l", false, "执行数据导入")
	listTable := flag.Bool("list", false, "查询数据库表")
	showTable := flag.String("show", "", "显示数据库表的 SQL 语句")
	reloadTable := flag.String("reload", "", "重新导入参数表")
	flag.Parse()
	if *querySQL != "" {
		sqlite3.Println(*querySQL)
	}
	if *execSQL != "" {
		err := sqlite3.ExecTx(sqlite3.NewTr(*execSQL))
		if err != nil {
			fmt.Println(err)
		}
	}
	if *load {
		Load()
	}
	if *listTable {
		sqlite3.Println("select name from sqlite_master where type=?", "table")
	}
	if *showTable != "" {
		sqlite3.Println("select sql from sqlite_master where type=? and name=?", "table", showTable)
	}
	if *reloadTable != "" {
		sqlite3.ExecTx(
			sqlite3.NewTr("delete from loadfile where name=?", reloadTable),
		)
	}
}
