package main

import (
	"database/sql"
	"fmt"
	"unsafe"

	// 引入 sqlite3 库
	_ "github.com/mattn/go-sqlite3"
)

// DB 数据库
type DB struct {
	sql.DB
}

// Open 打开数据库
func Open() (db *DB, err error) {
	var d *sql.DB
	d, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		return
	}
	db = (*DB)(unsafe.Pointer(d))
	return
}

// Fetch 查询，并读取结果
func (db *DB) Fetch(query string, args ...interface{}) {
	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var s string
	for rows.Next() {
		rows.Scan(&s)
		fmt.Println(s)
	}
}

func main() {
	db, err := Open()
	defer db.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	db.Fetch("select 'hello' ")
}
