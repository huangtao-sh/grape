package sqlite

import (
	"database/sql"
	"grape/path"
	"unsafe"

	// 引入 sqlite3 库
	_ "github.com/mattn/go-sqlite3"
)

// DB 数据库
type DB struct {
	sql.DB
}

var dataSourceName string

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

// Open 打开数据库，使用 Config 配置的文件
func Open() (db *DB, err error) {
	var d *sql.DB
	d, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return
	}
	db = (*DB)(unsafe.Pointer(d))
	return
}


