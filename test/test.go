package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"grape/data"
	"grape/gbk"
	"grape/path"
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

func convCdjy(s []string) (d []string) {
	if s[1] == "8" {
		d = append(d, s[0])
	} else {
		d = nil
	}
	return
}
func main() {
	file := path.NewPath(`C:\Users\huangtao\OneDrive\工作\参数备份\运营参数2020-02\YUNGUAN_MONTH_STG_TELLER_TRANSCONTROLS.del`)
	r, _ := file.Open()
	reader := gbk.NewReader(r)
	defer r.Close()
	re := text.NewReader(reader, false, text.NewSepSpliter(","), convCdjy)
	d := data.NewData()
	d.Add(1)
	go re.ReadAll(d)
	go func(d *data.Data) {
		defer d.Done()
		for row := range d.ReadCh() {
			fmt.Println(len(row), row)
		}
	}(d)
	d.Wait()
}
