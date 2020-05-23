package params

import (
	"bufio"
	"fmt"
	"grape/data"
	"grape/gbk"
	"grape/sqlite3"
	"grape/text"
	"grape/util"
	"io"
	"os"
	"strings"
	"sync"
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

// File tar、zip 压缩包获取文
type File interface {
	FileInfo() os.FileInfo
	Open() (io.ReadCloser, error)
}

// Loader 数据导入结构
type Loader struct {
	file           File
	name, ver, sql string
	data           *data.Data
}

// NewLoader 构造函数
func NewLoader(file File, name string, ver string, sql string) *Loader {
	return &Loader{file, name, ver, sql, data.NewData()}
}

// Check 检查文件是否已导入
func (l *Loader) Check(tx *sqlite3.Tx) (err error) {
	var count int
	info := l.file.FileInfo()
	err = tx.QueryRow("select count(name) from LoadFile where name=? and path=? and mtime>=datetime(?)",
		l.name, info.Name(), info.ModTime()).Scan(&count)
	util.CheckFatal(err)
	if count > 0 {
		return fmt.Errorf("文件 %s 已导入", info.Name())
	}
	return nil
}

// Save 导入后保存结果
func (l *Loader) Save(tx *sqlite3.Tx) (err error) {
	info := l.file.FileInfo()
	_, err = tx.Exec("insert or replace into LoadFile values(?,?,datetime(?),?)",
		l.name, info.Name(), info.ModTime(), l.ver)
	return
}

// Read 读取数据，并发送到通道中
func (l *Loader) Read() {
	defer l.data.Close()
	f, err := l.file.Open()
	util.CheckFatal(err)
	defer f.Close()
	r := gbk.NewReader(f)
	scanner := bufio.NewScanner(r)
	var fields []string
	for scanner.Scan() {
		s := scanner.Text()
		fields = strings.Split(s, ",")
		l.data.Write(text.Slice(fields)...)
	}
}

// Exec 执行导入操作
func (l *Loader) Exec(tx *sqlite3.Tx) (err error) {
	err = l.Check(tx)
	if err != nil {
		return
	}
	tx.Exec(fmt.Sprintf("delete from %s", l.name))
	l.data.Add(1)
	go l.data.Exec(tx, l.sql)
	go l.Read()
	l.data.Wait()
	l.Save(tx)
	return
}

// Load 导入数据
func (l *Loader) Load(wg *sync.WaitGroup) (err error) {
	defer wg.Done()
	err = sqlite3.ExecTx(l)
	if err == nil {
		fmt.Printf("%s 导入完成\n", l.name)
	} else {
		fmt.Println(err)
	}
	return
}

// Test 测试数据
func (l *Loader) Test() {
	l.data.Add(1)
	go l.data.Println()
	go l.Read()
	l.data.Wait()
}

// LoadCheck 检查文件是否已导入数据库
func LoadCheck(tx *sqlite3.Tx, name string, path File, ver string) (err error) {
	var info os.FileInfo
	var count int
	info = path.FileInfo()
	filename := info.Name()
	mtime := info.ModTime()
	err = tx.QueryRow("select count(name) from LoadFile where name=? and path=? and mtime>=datetime(?)", name, filename, mtime).Scan(&count)
	util.CheckFatal(err)
	if count > 0 {
		return fmt.Errorf("文件 %s 已导入", filename)
	}
	tx.Exec("insert or replace into LoadFile values(?,?,datetime(?),?)", name, filename, mtime, ver)
	return
}

// Checker 检查文件是否重复导入
type Checker struct {
	name string
	path File
	ver  string
}

// NewChecker 构造函数
func NewChecker(name string, path File, ver string) *Checker {
	return &Checker{name, path, ver}
}

// Exec 执行 SQL 语句
func (c *Checker) Exec(tx *sqlite3.Tx) error {
	return LoadCheck(tx, c.name, c.path, c.ver)
}

// GetVer 获取数据版本
func GetVer(name string) (ver string) {
	err := sqlite3.QueryRow("select ver from loadfile where name=?", name).Scan(&ver)
	util.CheckFatal(err)
	return
}

// PrintVer 打印数据版本
func PrintVer(name string) {
	fmt.Printf("数据版本：%s\n", GetVer(name))
}
