package loader

import (
	"fmt"
	"grape/data"
	"grape/sqlite3"
	"grape/text"
	"grape/util"
	"io"
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
	const (
		checkSQL = "select count(name) from LoadFile where name=? and path=? and mtime>=datetime(?)"
		doneSQL  = "insert or replace into LoadFile values(?,?,datetime(?),?)"
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

// File 文件接口
type File interface {
	FileInfo() os.FileInfo        // 获取文件信息
	Open() (io.ReadCloser, error) // 打开文件
}

// Reader 读取文件接口
type Reader interface {
	ReadAll(text.Data)
}

// NewReader Reader 构造函数
type NewReader func(r io.Reader) Reader

// Loader 数据装入类
type Loader struct {
	name, Ver, loadSQL string
	file               File
	new                NewReader
}

// NewLoader Loader 构造函数
func NewLoader(name string, ver string, loadSQL string, file File, new NewReader) *Loader {
	return &Loader{name, ver, loadSQL, file, new}
}

// ReadAll 读取所有数据
func (l *Loader) ReadAll(d text.Data) {
	r, err := l.file.Open()
	util.CheckFatal(err)
	defer r.Close()
	reader := l.new(r)
	reader.ReadAll(d)
}

// Test 导入主函数
func (l *Loader) Test() {
	d := data.NewData()
	d.Add(1)
	go l.ReadAll(d)
	go d.Println()
	d.Wait()
}

// Exec 执行导入操作
func (l *Loader) Exec(tx *sqlite3.Tx) (err error) {
	d := data.NewData()
	d.Add(1)
	go l.ReadAll(d)
	go tx.ExecCh(l.loadSQL, d)
	d.Wait()
	return
}

// Load 执行导入操作
func (l *Loader) Load() {
	info := l.file.FileInfo()
	err := sqlite3.ExecTx(
		loadCheck(l.name, info, l.Ver),
		sqlite3.NewTr(fmt.Sprintf("delete from %s", l.name)),
		l)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("导入文件 %s 完成！\n", info.Name())
	}
}
