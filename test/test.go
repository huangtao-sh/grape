package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"grape/sqlite3"
	"grape/text"
	"grape/util"
	"io"
	"os"
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

func main() {
	bf := `1,2,"3",4
"5","6",7,8
"9","10",11,12`
	b := bytes.NewReader([]byte(bf))
	r := text.NewReader(b, true, text.NewSepSpliter(","), text.UnQuote, text.Include(0, 1))
	for r.Next() {
		fmt.Println(r.Read()...)
	}

}
