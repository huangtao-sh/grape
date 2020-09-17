package sqlite3

import (
	"database/sql"
	"errors"
	"grape/util"
)

// Tx 数据库事务
type Tx struct {
	*sql.Tx
}

// NewTx 新建事务
func NewTx() *Tx {
	tx, err := NewDB().Begin()
	util.CheckFatal(err)
	return &Tx{tx}
}

// FetchValue 查询值
func (tx *Tx) FetchValue(sql string, args ...interface{}) interface{} {
	return fetchValue(tx, sql, args...)
}

// Fetch 查询数据
func (tx *Tx) Fetch(sql string, args ...interface{}) (reader *RowReader) {
	return fetch(tx, sql, args...)
}

// DataCh 数据通道接口
type DataCh interface {
	Done()
	ReadCh() <-chan []interface{}
}

// ExecCh 执行通道中的数据
func (tx *Tx) ExecCh(sql string, data DataCh) (err error) {
	defer data.Done()
	stmt, err := tx.Prepare(sql)
	util.CheckFatal(err)
	defer stmt.Close()
	for row := range data.ReadCh() {
		_, err = stmt.Exec(row...)
		if err != nil {
			return
		}
	}
	return
}

// Txer 数据库事务执行接口
type Txer interface {
	Exec(tx *Tx) error
}

// ExecFunc 执行函数类型
type ExecFunc func(*Tx) error

// ExecTx 执行事务
func ExecTx(txers ...interface{}) (err error) {
	tx := NewTx()
	defer tx.Rollback()
	for _, txer := range txers {
		switch execer := txer.(type) {
		case ExecFunc:
			err = execer(tx)
		case Txer:
			err = execer.Exec(tx)
		default:
			return errors.New("执行错误")
		}
		if err != nil {
			return
		}
	}
	tx.Commit()
	return
}

// NewTr Tr构造函数
func NewTr(sql string, params ...interface{}) ExecFunc {
	return func(tx *Tx) (err error) {
		_, err = tx.Exec(sql, params...)
		return
	}
}
