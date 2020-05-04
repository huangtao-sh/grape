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
