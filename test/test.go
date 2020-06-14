package main

import (
	"archive/tar"
	"compress/gzip"
	"grape/data"
	"grape/sqlite3"
	"grape/text"
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
	file := `C:\Users\huangtao\OneDrive\工作\参数备份\运营参数2020-02\YUNGUAN_MONTH_STG_ZSRUN_GGJGM.del`
	r, err := os.Open(file)
	util.CheckFatal(err)
	defer r.Close()
	reader := text.NewReader(text.Decode(r, false, true), false, text.NewSepSpliter(","),
		text.Include(0, 1, 3-43, 7-43, 15-43, 16-43, 17-43))
	d := data.NewData()
	d.Add(1)
	go d.Println()
	go reader.ReadAll(d)
	d.Wait()
}
