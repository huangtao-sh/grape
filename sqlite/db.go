package sqlite

import (
	"database/sql"

	// 引入 sqlite3 库
	_ "github.com/mattn/go-sqlite3"
)

// Open 打开数据库
func Open(dataSourceName string) (*sql.DB, error) {
	return sql.Open("sqlite3", dataSourceName)
}

// DB 数据库接口，可以是 sql.Tx 或 sql.DB
type DB interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Prepare(string) (*sql.Stmt, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
}

// DataProvider 提供数据，供 ExecMany 使用
type DataProvider interface {
	Next() bool
	Read() ([]interface{}, error)
}

// ExecMany 执行多条 sql 语句
func ExecMany(db DB, sql string, data DataProvider) (err error) {
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
