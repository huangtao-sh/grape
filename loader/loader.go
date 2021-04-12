package loader

import (
	"fmt"
	"grape/sqlite3"
	"grape/util"
	"io"
	"log"
	"os"
	"sync"
)

var initLoadfile = sync.Once{}

func createLoadFile() {
	sqlite3.ExecScript(`
create table if not exists loadfile(
	name	text 	primary key,  -- 类型
	path	text,	              -- 文件名
	mtime	text,                 -- 文件修改时间
	ver		text		          -- 文件版本
)`)
}

// LoadCheck 检查文件是否已导入数据库
func loadCheck(name string, info os.FileInfo, ver string) sqlite3.ExecFunc {
	const (
		checkSQL = "select count(name) from loadfile where name=? and path=? and mtime>=datetime(?)"
		doneSQL  = "insert or replace into loadfile values(?,?,datetime(?),?)"
	)
	initLoadfile.Do(createLoadFile) // 仅在调用 LoadCheck 时运行一次建表
	var count int
	filename := info.Name()
	mtime := info.ModTime()
	return func(tx *sqlite3.Tx) (err error) {
		err = tx.QueryRow(checkSQL, name, filename, mtime).Scan(&count)
		util.CheckFatal(err)
		if count > 0 {
			return fmt.Errorf("文件 %s 已导入", filename)
		}
		tx.Exec(doneSQL, name, filename, mtime, ver)
		return
	}
}

// Loader 数据装入类
type Loader struct {
	name, sql string      // 表名，导入的SQL语句
	info      os.FileInfo // 文件信息
	Ver       string      // 导入数据版本
	data      Reader      // 数据读取程序
	Clear     bool        // 是否清理数据库，默认为是
	Check     bool        // 是否需要检查文件导入
}

// NewLoader Loader 构造函数
func NewLoader(info os.FileInfo, name string, sql string, data Reader) *Loader {
	return &Loader{name: name, info: info, sql: sql, data: data, Clear: true, Check: true}
}

// Test 导入主函数
func (l *Loader) Test() {
	var (
		columns []string
		err     error
	)
	for i := 0; i < 10 && err == nil; columns, err = l.data.Read() {
		if columns != nil {
			fmt.Println(Slice(columns)...)
			i++
		}
	}
	if err != nil && err != io.EOF {
		fmt.Println(err)
	}
}

// Exec 执行导入操作
func (l *Loader) Exec(tx *sqlite3.Tx) (err error) {
	var columns []string
	stmt, err := tx.Prepare(l.sql)
	log.Printf("准备SQL：%s\n", l.sql)
	if err != nil {
		return
	}
	defer stmt.Close()
	for ; err == nil; columns, err = l.data.Read() {
		if columns != nil {
			_, err = stmt.Exec(Slice(columns)...)
			if err != nil {
				return
			}
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

// Load 导入数据
func (l *Loader) Load() error {
	var steps []interface{}
	if l.Check {
		steps = append(steps, loadCheck(l.name, l.info, l.Ver))
	}
	if l.Clear {
		steps = append(steps, sqlite3.NewTr(fmt.Sprintf("delete from %s", l.name)))
	}
	steps = append(steps, l)
	return sqlite3.ExecTx(steps...)
}
