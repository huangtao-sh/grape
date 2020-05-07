package params

import (
	"errors"
	"fmt"
	"grape/sqlite3"
	"os"
)

func init() {
	sqlite3.Config(":memory:")
	sqlite3.ExecScript(`
	create table if not exists LoadFile(
		name	text 	primary key,  -- 类型
		path	text,	              -- 文件名
		mtime	text,             -- 文件修改时间
		ver		text		          -- 文件版本
	)
	`)
}

// Filer tar、zip 压缩包获取文
type Filer interface {
	FileInfo() os.FileInfo
}

// LoadCheck 检查文件是否已导入数据库
func LoadCheck(tx *sqlite3.Tx, name string, path interface{}, ver string) (err error) {
	var info os.FileInfo
	var value string
	switch file := path.(type) {
	case string:
		info, err = os.Lstat(file)
		if err != nil {
			return
		}
	case Filer:
		info = file.FileInfo()
	}
	filename := info.Name()
	mtime := info.ModTime().Format("2006-01-02 15:04:05") // 把时间格式化成字符串
	fmt.Println(filename, mtime, ver)
	row := tx.QueryRow("select name from LoadFile where name=? and path=? and mtime>=?", name, path, mtime)
	err = row.Scan(&value)
	fmt.Println(err)
	if err == nil {
		return errors.New("文件已导入")
	}
	tx.Exec("insert or replace into LoadFile values(?,?,?,?)", name, filename, mtime, ver)
	return nil
}
