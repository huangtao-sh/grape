package sqlite

import (
	"database/sql"
	"grape/util"
)

// Execer 数据库接口，可以是 sql.Tx 或 sql.DB
type Execer interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Prepare(string) (*sql.Stmt, error)
}

//Querier 数据库接口，可以为sql.Tx 或 sql.DB
type Querier interface {
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
}

// Dber 数据接口
type Dber interface {
	Execer
	Querier
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
	defer stmt.Close()
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

// ExecCh 从通道中读取数据并执行
func ExecCh(db Execer, sql string, data util.Data) {
	stmt, err := db.Prepare(sql)
	util.CheckFatal(err)
	for row := range data.Read() {
		_, err := stmt.Exec(row...)
		util.CheckFatal(err)
	}
}

// FetchCh 执行查询，返回多条记录
func FetchCh(db Querier, data *util.Data, query string, args ...interface{}) {
	defer data.Close()                    // 退出时关闭数据通道
	ch := data.Write()                    // 获取写入通道
	rows, err := db.Query(query, args...) // 执行查询
	util.CheckFatal(err)
	columns, err := rows.Columns()
	util.CheckFatal(err)
	count := len(columns)
	values := make([]interface{}, count)
	addrs := make([]interface{}, count)
	for i := range columns {
		addrs[i] = &values[i]
	}
	for rows.Next() {
		rows.Scan(addrs...) // 读取数据
		ch <- values        // 发送数据
	}
}

// FetchValue 执行查询，并返回一个值
func FetchValue(db Querier, query string, args ...interface{}) (value interface{}) {
	row := db.QueryRow(query, args...)
	row.Scan(&value)
	return
}
