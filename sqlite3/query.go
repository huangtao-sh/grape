package sqlite3

import (
	"database/sql"
	"reflect"
)


// 查询接口
type rowReader struct {
	rows          *sql.Rows
	addrs, values []interface{}
}

// 查询一条记录
func (db DB) FindOne(query string, args ...interface{}) *sql.Row {
	return db.db.QueryRow(query, args...)
}

// 查询值
func (db DB) FindValue(query string, args ...interface{}) (result interface{}) {
	row := db.db.QueryRow(query, args...)
	row.Scan(&result)
	return
}

// 执行查询，返回多条记录
func (db DB) Fetch(query string, args ...interface{}) (reader *rowReader, err error) {
	rows, err := db.db.Query(query, args...)
	if err != nil {
		return
	}
	columns, err := rows.Columns()
	if err != nil {
		return
	}
	count := len(columns)
	values := make([]interface{}, count)
	addrs := make([]interface{}, count)
	for i := 0; i < len(columns); i++ {
		addrs[i] = &values[i]
	}
	reader = &rowReader{rows, addrs, values}
	return
}

// 执行查询，返回多条记录
func (db DB) FetchValue(query string,addr interface{},args ...interface{}) (err error) {
	row:= db.db.QueryRow(query, args...)
	err=row.Scan(addr)
	return
}

// 是否有下一条记录
func (reader *rowReader) Next() bool {
	return reader.rows.Next()
}

// 关闭查询
func (reader *rowReader) Close() error {
	return reader.rows.Close()
}

// 读取下一条记录
func (reader *rowReader) Read() []interface{} {
	reader.rows.Scan(reader.addrs...)
	return reader.values
}

// 读取下一条记录
func (reader *rowReader) Scan(addr ...interface{}) error {
	return reader.rows.Scan(addr...)
}

// 获取 Struct 的地址列表
func StructAddr(s interface{})[]interface{}{
	k:=reflect.ValueOf(s).Elem()
	result:=make([]interface{},k.NumField())
	for i:=0;i<k.NumField();i++{
		result[i]=k.Field(i).Addr().Interface()
	}
	return result
}