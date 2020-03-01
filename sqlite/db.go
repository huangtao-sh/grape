package sqlite

import (
	"database/sql"
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
func Config(path string) {
	dataSourceName = path
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

// Execer 数据库接口，可以是 sql.Tx 或 sql.DB
type Execer interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Prepare(string) (*sql.Stmt, error)
}

// DataProvider 提供数据，供 ExecMany 使用
type DataProvider interface {
	Next() bool
	Read() ([]interface{}, error)
}

// ExecMany 执行多条 sql 语句
func ExecMany(db Execer, sql string, data DataProvider) (err error) {
	var row []interface{}
	stmt, err := db.Prepare(sql)
	if err != nil {
		return
	}
	for data.Next() {
		row, err = data.Read()
		if err != nil {
			return
		}
		_, err = stmt.Exec(row)
		if err != nil {
			return
		}
	}
	return
}
