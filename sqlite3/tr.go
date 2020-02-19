package sqlite3

import "database/sql"

// 单条语句命令
type execOne struct {
	sql  string
	args []interface{}
}

// 单条语句执行
func (one execOne) exec(tx *sql.Tx) (err error) {
	_, err = tx.Exec(one.sql, one.args...)
	return
}

// 多条语句命令
type execMany struct {
	sql string
	r   reader
}

// 多条语句执行
func (many execMany) exec(tx *sql.Tx) (err error) {
	stmt, err := tx.Prepare(many.sql) // 执行准备语句
	if err != nil {
		return
	}
	defer stmt.Close() // 清理准备语句
	data := many.r
	for {
		values, err := data.Read() // 切换到下一行，并读取数据
		if err == nil {
			_, err = stmt.Exec(values...) // 读取成功则执行
		} else if err.Error() == "EOF" {
			return nil //数据读取完毕，正常返回
		}
		if err != nil {
			return err //遇到错误，则返回错误
		}
	}
}

// 执行命令接口
type execer interface {
	exec(tx *sql.Tx) error
}

// 事务接口
type Tr struct {
	db    DB
	tasks []execer
}

// 进入事务
func (db DB) Tran() (tr *Tr) {
	return &Tr{db, []execer{}}
}

// 当前事务中加入单条语句
func (tr *Tr) Add(sql string, args ...interface{}) {
	tr.tasks = append(tr.tasks, execOne{sql, args})
}

// 当前事务中加入多条语句
func (tr *Tr) AddMany(sql string, r reader) {
	tr.tasks = append(tr.tasks, execMany{sql, r})
}

// 执行完整事务
func (tr *Tr) Exec() (err error) {
	tx, err := tr.db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()
	for _, execer := range tr.tasks {
		err = execer.exec(tx)
		if err != nil {
			return
		}
	}
	tx.Commit()
	return
}

// rader 接口
type reader interface {
	Read() ([]interface{}, error)
}
