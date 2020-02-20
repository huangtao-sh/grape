package main

import (
	"flag"
	"fmt"
	"grape/xyk"
)

var auto, init_, load, compare, export *bool
var sql *string

func init() {
	auto = flag.Bool("a", false, "自动执行相关操作")
	init_ = flag.Bool("i", false, "初始化数据库")
	load = flag.Bool("l", false, "导入数据")
	compare = flag.Bool("c", false, "执行对账程序")
	export = flag.Bool("e", false, "导出报表")
	sql = flag.String("s", "", "执行 sql 语句")
	flag.Parse()
}
func main() {
	if *init_ {
		xyk.CreateDB()
		fmt.Println("创建数据库表完成！")
	}
	if *auto || *load {
		xyk.Load()
	}
	if *sql != "" {
		xyk.Query(*sql)
	}

}
