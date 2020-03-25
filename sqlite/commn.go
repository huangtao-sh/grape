package sqlite

import "database/sql"

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
