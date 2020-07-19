package main

import (
	"archive/tar"
	"compress/gzip"
	_ "grape/params"
	"grape/params/jym"
	"grape/path"
	"grape/sqlite3"
	"grape/util"
	"io"
	"os"
)

// MMain Test
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

func main() {
	p := path.NewPath(`C:\Users\huangtao\OneDrive\工作\参数备份\生产参数\交易码参数.xlsx`)
	loader := jym.LoadJycs(p, "1.0")
	loader.Load()
	sqlite3.Println("select *,rowid from jymcs")

	/*
		sqlite3.Config(":memory:")
		sqlite3.ExecScript(`create table test(a text)`)
		err := sqlite3.ExecTx(
			sqlite3.NewTr("insert or replace into test(a,rowid)values(?,?)",  "test",nil),
			sqlite3.NewTr("insert or replace into test(a,rowid)values(?,?)", "Ab", nil),
			sqlite3.NewTr("insert or replace into test(a,rowid)values(?,?)",  "hello","2"),
			//sqlite3.NewTr("insert into test(rowid,a)values(?,?)", "", "cd"),
		)
		if err != nil {
			fmt.Println(err)
		}
		sqlite3.Println("select rowid,* from test")
	*/
}
