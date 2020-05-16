package sqlite3

import (
	"database/sql"
	"grape/data"
	"grape/util"
	"unsafe"
)

// Tx 数据库事务
type Tx struct {
	sql.Tx
}

// NewTx 新建事务
func NewTx() *Tx {
	tx, err := NewDB().Begin()
	util.CheckFatal(err)
	return (*Tx)(unsafe.Pointer(tx))
}

// ExecCh 批量执行 sql 语句，数据从通道中获取
func (tx *Tx) ExecCh(sql string, d *data.Data) {
	defer d.Done()
	stmt, err := tx.Prepare(sql)
	util.CheckFatal(err)
	defer stmt.Close()
	for row := range d.ReadCh() {
		_, err = stmt.Exec(row...)
		util.CheckFatal(err)
	}
}

// FetchValue 查询值
func (tx *Tx) FetchValue(sql string, args ...interface{}) interface{} {
	return fetchValue(tx, sql, args...)
}

// Txer 数据库事务执行接口
type Txer interface {
	Exec(tx *Tx) error
}

// ExecTx 执行事务
func ExecTx(txers ...Txer) (err error) {
	tx := NewTx()
	defer tx.Rollback()
	for _, txer := range txers {
		err = txer.Exec(tx)
		if err != nil {
			return
		}
	}
	tx.Commit()
	return
}

// Tr 在事务中执行语句
type Tr struct {
	sql    string
	params []interface{}
}

// NewTr Tr构造函数
func NewTr(sql string, params ...interface{}) *Tr {
	return &Tr{sql, params}
}

// Exec 执行 Tr
func (t *Tr) Exec(tx *Tx) (err error) {
	_, err = tx.Exec(t.sql, t.params...)
	return
}
