package sqlite3

import (
	"database/sql"
	"fmt"
	"grape/util"
	"reflect"
)

// RowReader 查询接口
type RowReader struct {
	rows          *sql.Rows
	addrs, values []interface{}
}

// Next 是否有下一条记录
func (reader *RowReader) Next() bool {
	return reader.rows.Next()
}

// Close 关闭查询
func (reader *RowReader) Close() error {
	return reader.rows.Close()
}

// Read 读取下一条记录
func (reader *RowReader) Read() []interface{} {
	reader.rows.Scan(reader.addrs...)
	return reader.values
}

// Scan 读取下一条记录
func (reader *RowReader) Scan(addr ...interface{}) error {
	return reader.rows.Scan(addr...)
}

type querier interface {
	QueryRow(string, ...interface{}) *sql.Row
	Query(string, ...interface{}) (*sql.Rows, error)
}

// fetch 执行查询，并返回查询结果
func fetch(db querier, sql string, args ...interface{}) (reader *RowReader) {
	rows, err := db.Query(sql, args...)
	util.CheckFatal(err)
	columns, err := rows.Columns()
	util.CheckFatal(err)
	count := len(columns)
	values := make([]interface{}, count)
	addrs := make([]interface{}, count)
	for i := 0; i < len(columns); i++ {
		addrs[i] = &values[i]
	}
	reader = &RowReader{rows, addrs, values}
	return
}

// fetchValue 执行查询，返回多条记录
func fetchValue(db querier, query string, args ...interface{}) (value interface{}) {
	row := db.QueryRow(query, args...)
	err := row.Scan(&value)
	util.CheckFatal(err)
	return
}

// Fetch 执行查询，并返回查询结果
func Fetch(sql string, args ...interface{}) (reader *RowReader) {
	return fetch(NewDB(), sql, args...)
}

// FetchValue 执行查询，并返回值
func FetchValue(sql string, args ...interface{}) interface{} {
	return fetchValue(NewDB(), sql, args...)
}

// Println 执行查询，并打印查询结果
func Println(sql string, args ...interface{}) {
	rows := Fetch(sql, args...)
	for rows.Next() {
		fmt.Println(rows.Read()...)
	}
}

// Printf 执行查询，并打印查询结果
func Printf(format string, sql string, args ...interface{}) {
	rows := Fetch(sql, args...)
	for rows.Next() {
		fmt.Printf(format, rows.Read()...)
	}
}

// StructAddr 获取 Struct 的地址列表
func StructAddr(s interface{}) []interface{} {
	k := reflect.ValueOf(s).Elem()
	result := make([]interface{}, k.NumField())
	for i := 0; i < k.NumField(); i++ {
		result[i] = k.Field(i).Addr().Interface()
	}
	return result
}
