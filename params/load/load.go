package load

import (
	"fmt"
	"grape/data"
	"grape/sqlite3"
	"grape/text"
	"grape/util"
	"os"
	"sync"
)

var initLoadfile = sync.Once{}

func createLoadFile() {
	sqlite3.ExecScript(`
create table if not exists LoadFile(
	name	text 	primary key,  -- 类型
	path	text,	              -- 文件名
	mtime	text,                 -- 文件修改时间
	ver		text		          -- 文件版本
)`)
}

// LoadCheck 检查文件是否已导入数据库
func loadCheck(name string, info os.FileInfo, ver string) sqlite3.ExecFunc {
	initLoadfile.Do(createLoadFile) // 仅在调用 LoadCheck 时运行一次建表
	var count int
	filename := info.Name()
	mtime := info.ModTime()
	return func(tx *sqlite3.Tx) (err error) {
		err = tx.QueryRow("select count(name) from LoadFile where name=? and path=? and mtime>=datetime(?)", name, filename, mtime).Scan(&count)
		util.CheckFatal(err)
		if count > 0 {
			return fmt.Errorf("文件 %s 已导入", filename)
		}
		tx.Exec("insert or replace into LoadFile values(?,?,datetime(?),?)", name, filename, mtime, ver)
		return
	}
}

// Reader 读取数据
type Reader interface {
	ReadAll(text.Data)
}

// Loader 数据导入模型
type Loader struct {
	Name             string
	FileInfo         os.FileInfo
	ver              string
	reader           Reader
	initSQL, loadSQL string
}

// NewLoader 构造函数
func NewLoader(name string, fileInfo os.FileInfo, ver string, txt Reader, initSQL string, loadSQL string) *Loader {
	return &Loader{name, fileInfo, ver, txt, initSQL, loadSQL}
}

// Load 导入数据
func (l *Loader) Load() {
	if l.initSQL != "" {
		sqlite3.ExecScript(l.initSQL) // 初始化数据库
	}
	err := sqlite3.ExecTx(
		loadCheck(l.Name, l.FileInfo, l.ver),
		sqlite3.NewTr(fmt.Sprintf("delete from %s", l.Name)),
		l)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("导入文件 %s 完成！\n", l.FileInfo.Name())
	}
}

// Exec 执行导入操作
func (l *Loader) Exec(tx *sqlite3.Tx) (err error) {
	d := data.NewData()
	d.Add(1)
	go tx.ExecCh(l.loadSQL, d)
	go l.reader.ReadAll(d)
	d.Wait()
	return
}

// Test 执行测试
func (l *Loader) Test() {
	d := data.NewData()
	d.Add(1)
	go l.reader.ReadAll(d)
	go d.Println()
	d.Wait()
}
