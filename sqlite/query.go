package sqlite

import (
	"database/sql"
	"fmt"
)

// Reader 查询数据行
type Reader struct {
	Rows          *sql.Rows
	addrs, values []interface{}
}

// Next 是否有下一条记录
func (r *Reader) Next() bool {
	return r.Rows.Next()
}

// Close 关闭查询
func (r *Reader) Close() error {
	return r.Rows.Close()
}

// Read 读取下一条记录
func (r *Reader) Read() []interface{} {
	r.Rows.Scan(r.addrs...)
	return r.values
}

// Scan 读取下一条记录
func (r *Reader) Scan(addr ...interface{}) error {
	return r.Rows.Scan(addr...)
}

// Query 执行查询，返回多条记录
func Query(db DB, query string, args ...interface{}) (reader *Reader, err error) {
	rows, err := db.Query(query, args...)
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
	for i := range columns {
		addrs[i] = &values[i]
	}
	reader = &Reader{rows, addrs, values}
	return
}

// ExecQuery 执行查询，并输出
func ExecQuery(db DB, sql string, args ...interface{}) {
	r, err := Query(db, sql, args...)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer r.Close()
	for r.Next() {
		fmt.Println(r.Read()...)
	}
	return
}
