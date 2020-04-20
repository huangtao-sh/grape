package sqlite3

import (
	"database/sql"
	"grape/path"
	"grape/util"

	// 引用 go-sqlite3 进行初始化
	_ "github.com/mattn/go-sqlite3"
)

var dataSourceName string
var db *sql.DB

// Config  配置数据库文件
// 配置完成之后，可以直接调用 Open 打开
func Config(pathName string) {
	if (pathName != ":memory:") && (path.NewPath(pathName).Dir() == ".") { // path 如果不是 :memory:，无目录的指定默认目录
		dataHome := path.Home.Join(".data")
		dataHome.Ensure() // 目录不存在则自动创建
		pathName = (dataHome.Join(pathName)).String()
	}
	dataSourceName = pathName
}

// NewDB 打开数据库
func NewDB() *sql.DB {
	if db == nil {
		var err error
		db, err = sql.Open("sqlite3", dataSourceName)
		util.CheckFatal(err)
	}
	return db
}

// Close 关闭数据库
func Close() {
	if db != nil {
		db.Close()
		db = nil
	}
}

// ExecScripts 执行脚本，支持多条语句，用于执行 DDL 语句
func ExecScripts(sql string) {
	_, err := NewDB().Exec(sql)
	util.CheckFatal(err)
}
