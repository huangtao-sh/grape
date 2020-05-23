package main

import (
	"flag"
	"fmt"
	_ "grape/params"
	"grape/sqlite3"
)

func main() {
	execSQL := flag.String("e", "", "执行 SQL 语句")
	querySQL := flag.String("q", "", "执行查询")
	load := flag.Bool("l", false, "执行数据导入")
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
}
