package sqlite3

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// 数据库
type DB struct {
	db *sql.DB
}

// 打开数据库
func Open(dataSourceName string) (db DB, err error) {
	db_, err := sql.Open("sqlite3", dataSourceName)
	db = DB{db_}
	return
}

// 关闭数据库
func (db DB) Close() error {
	return db.db.Close()
}

// 开启事务
func (db DB) Begin() (*sql.Tx, error) {
	return db.db.Begin()
}

// 执行 SQL 语句
func (db DB) Exec(statment string, args ...interface{}) (sql.Result, error) {
	return db.db.Exec(statment, args...)
}
