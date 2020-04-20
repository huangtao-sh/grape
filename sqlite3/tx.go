package sqlite3

import (
	"database/sql"
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
