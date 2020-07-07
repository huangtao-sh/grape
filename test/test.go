package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	_ "grape/params"
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
	fmt.Println(util.Match(`\d{3}`, "134254325"))
	fmt.Println(util.FullMatch(`\d{3}`, "134254325"))
	fmt.Println(util.FullMatch(`\w{3}`, "ab1"))
}
