package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"grape/params/km"
	"grape/path"
	"grape/sqlite"
	"grape/sqlite3"
	"grape/util"
	"io"
	"os"
	"sync"
	"time"
)

func MMain() {
	sqlite3.Config(":memory:")
	var t *tar.Reader
	sqlite3.ExecScript("create table if not exists test(name text,modtime text)")
	f, err := os.Open("E:/yg20200430.tar.gz")
	util.CheckFatal(err)
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		f.Seek(0, 0)
		t = tar.NewReader(f)
	} else {
		t = tar.NewReader(gz)
	}
	func() {
		tx := sqlite3.NewTx()
		defer tx.Rollback()
		stmt, err := tx.Prepare("insert into test values(?,datetime(?))")
		util.CheckFatal(err)
		defer stmt.Close()
		for h, err := t.Next(); err != io.EOF; h, err = t.Next() {
			info := h.FileInfo()
			stmt.Exec(info.Name(), info.ModTime())
		}
		tx.Commit()
	}()
	sqlite3.Println("select modtime,name from test")
}

func Zip() {
	sqlite.Config(":memory:")
	db, _ := sqlite.Open()
	defer db.Close()
	db.Exec("create table if not exists test(name text,modtime text)")
	f, err := zip.OpenReader("E:/OneDrive/工作/参数备份/运营参数2020-03.zip")
	util.CheckFatal(err)
	defer f.Close()
	stmt, err := db.Prepare("insert into test values(?,datetime(?))")
	for _, file := range f.File {
		info := file.FileInfo()
		stmt.Exec(info.Name(), info.ModTime())
	}
	db.ExecQuery("select modtime,name from test")
}

func insert(db *sqlite.DB, t time.Time) {
	tx, err := db.Begin()
	util.CheckFatal(err)
	defer tx.Rollback()
	tx.Exec("insert into test values(?,?)", "Hello", t)
	tx.Commit()
}
func Load() {
	sqlite.Config(":memory:")
	db, err := sqlite.Open()
	util.CheckFatal(err)
	defer db.Close()
	db.Exec(`
	create table if not exists test(
		name text	primary key,
		mtime	datetime
	)
	`)
	t := time.Now()
	insert(db, t)
	row := db.QueryRow("select name from test where datetime(mtime)>=datetime(?,'1 minutes')", t)
	var s string
	err = row.Scan(&s)
	fmt.Println(err, s, s == "")
}

func main() {
	//sqlite3.Config(":memory:")
	//file := path.NewPath("E:/OneDrive/工作/参数备份/交易菜单").Find("menu*.xml")
	file := path.NewPath(path.NewPath("E:/OneDrive/工作/参数备份/运营参数2020-04").Find("YUNGUAN_MONTH_STG_ZSRUN_GGNBZHMB.*"))
	fmt.Println(file)
	loader := km.NewNbzhmbLoader(file)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	err := loader.Load(wg)
	if err != nil {
		fmt.Println(err)
	}
	wg.Wait()
	sqlite3.Println("select * from nbzhmb where km=? and xh=?", "511024", 1)
}
