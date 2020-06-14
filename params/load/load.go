package load

import (
	"fmt"
	"grape/data"
	"grape/sqlite3"
	"grape/text"
	"grape/util"
	"io"
	"os"
)

func init() {
	sqlite3.Config("params.db")
	sqlite3.ExecScript(`
	create table if not exists LoadFile(
		name	text 	primary key,  -- 类型
		path	text,	              -- 文件名
		mtime	text,                 -- 文件修改时间
		ver		text		          -- 文件版本
	)
	`)
}

// LoadCheck 检查文件是否已导入数据库
func LoadCheck(name string, path File, ver string) sqlite3.ExecFunc {
	var info os.FileInfo
	var count int
	info = path.FileInfo()
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

// File tar、zip 压缩包获取文
type File interface {
	FileInfo() os.FileInfo
	Open() (io.ReadCloser, error)
}

// Reader 读取数据
type Reader interface {
	ReadAll(text.Data)
}

// Loader 数据导入模型
type Loader struct {
	Name string
	File
	ver              string
	reader           Reader
	initSQL, loadSQL string
}

// NewLoader 构造函数
func NewLoader(name string, file File, ver string, txt *text.Reader, initSQL string, loadSQL string) *Loader {
	return &Loader{name, file, ver, txt, initSQL, loadSQL}
}

// Load 导入数据
func (l *Loader) Load() {
	sqlite3.ExecScript(l.initSQL) // 初始化数据库
	err := sqlite3.ExecTx(
		LoadCheck(l.Name, l.File, l.ver),
		sqlite3.NewTr(fmt.Sprintf("delete from %s", l.Name)),
		l)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("导入文件 %s 完成！\n", l.File.FileInfo().Name())
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
