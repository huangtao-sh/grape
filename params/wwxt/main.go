package wwxt

import (
	"flag"
	"fmt"
	"grape/path"
	"grape/sqlite3"
	"runtime"
)

// Root 参数根目录
var Root *path.Path

func init() {
	Root = path.NewPath("~/OneDrive/工作/参数备份/外围系统")
}

// Version 打印版本信息
func Version() {
	fmt.Println("外围系统参数程序 Ver 0.1")
	fmt.Printf("Compiled by %s\n", runtime.Version())
}

// Main 主程序
func Main() {
	defer sqlite3.Close() // 关闭数据库，释放资源
	load := flag.Bool("l", false, "导入数据")
	export := flag.Bool("e", false, "导出数据")
	sql := flag.String("E", "", "执行 SQL 语句")
	names := flag.String("a", "", "新增外围系统名称，如有多个系统用逗号分开")
	show := flag.Bool("s", false, "显示外围系统列表")
	query := flag.String("q", "", "查询指定外围系统")
	version := flag.Bool("V", false, "显示版本信息")
	flag.Parse()
	if *load {
		Load()
	}
	if *show {
		fmt.Println("编号  维护日期     机构范围    机构码       系统名称")
		sqlite3.Printf("%3d  %-10s  %-8s %-11s %s\n", "select id,date,jglx,jgh,name from wwxt order by id")
	}
	if *names != "" {
		Add(names)
	}
	if *sql != "" {
		sqlite3.Println(*sql)
	}
	if *query != "" {
		fmt.Println("编号  维护日期     机构范围    机构码       系统名称")
		sqlite3.Printf("%3d  %-10s  %-8s %-11s %s\n", "select id,date,jglx,jgh,name from wwxt where name like ? order by id",
			fmt.Sprintf("%%%s%%", *query))
	}
	if *export {
		Export()
	}
	if *version {
		Version()
	}
}
