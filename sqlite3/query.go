package sqlite3

import (
	"database/sql"
	"fmt"
	"grape/util"
	"reflect"
	"strings"
)

// RowReader 查询接口
type RowReader struct {
	*sql.Rows
	addrs, values []interface{}
}

// Read 读取下一条记录
func (r *RowReader) Read() []interface{} {
	r.Scan(r.addrs...)
	return r.values
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

// QueryRow 执行查询，返回一行数据
func QueryRow(query string, args ...interface{}) *sql.Row {
	return NewDB().QueryRow(query, args...)
}

// PrintRow 查询并打印单行数据，作为一个对象打印
// 采用的格式为： Name     Value 的格式打印
func PrintRow(header string, query string, args ...interface{}) (err error) {
	var width int
	headers := strings.Split(header, ",")
	count := len(headers)
	values := make([]interface{}, count)
	addrs := make([]interface{}, count)
	for i := range headers {
		addrs[i] = &values[i]
		if l := util.Wlen(headers[i]); l > width {
			width = l
		}
	}
	format := fmt.Sprintf("%%%ds  %%s", width)
	err = QueryRow(query, args...).Scan(addrs...)
	if err != nil {
		return
	}
	for i, header := range headers {
		fmt.Println(util.Sprintf(format, header, fmt.Sprint(values[i])))
	}
	return
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
		//fmt.Printf(format, rows.Read()...)
		fmt.Print(util.Sprintf(format, rows.Read()...))
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
