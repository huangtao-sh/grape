package params

import (
	"fmt"
	"grape/sqlite3"
	"grape/util"
)

func init() {
	sqlite3.Config("params.db")
}

// GetVer 获取数据版本
func GetVer(name string) (ver string) {
	err := sqlite3.QueryRow("select ver from loadfile where name=?", name).Scan(&ver)
	util.CheckFatal(err)
	return
}

// PrintVer 打印数据版本
func PrintVer(name string) {
	fmt.Printf("数据版本：%s\n", GetVer(name))
}
