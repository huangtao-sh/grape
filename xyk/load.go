package xyk

import (
	"archive/zip"
	"compress/gzip"
	"fmt"
	"grape/gbk"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

var Home string //程序根目录

func init() {
	home, _ := os.UserHomeDir()
	Home = filepath.Join(home, "信用卡对账")
}

type Reader struct {
	reader   *zip.ReadCloser
	date     string      // 日期
	trac     *zip.File   // 银数日报表
	jorj     *zip.File   // 银数会计流水
	eve      *zip.File   // 银数交易日志
	rd1002   *zip.File   // 银联日报表
	indfiles []*zip.File // 银联明细文件
}

func ReadAll(file *zip.File) (data []byte, err error) {
	var reader io.Reader
	f, err := file.Open()
	if err != nil {
		return
	}
	defer f.Close()
	if filepath.Ext(file.Name) == ".gz" {
		reader, _ = gzip.NewReader(f)
	} else {
		reader = f
	}
	reader = gbk.NewReader(reader)
	return ioutil.ReadAll(reader)
}

func OpenReader(path string) (reader *Reader, err error) {
	date := "20" + filepath.Base(path)[:6]
	zreader, err := zip.OpenReader(path)
	if err != nil {
		return
	}
	reader = &Reader{reader: zreader, date: date}
	for _, path := range zreader.File {
		name := filepath.Base(path.Name)
		switch {
		case name == "0316-EVE-"+date:
			reader.eve = path
		case len(name) >= 13 && name[:13] == "GLREPORT-JORJ":
			reader.jorj = path
		case len(name) >= 13 && name[:13] == "GLREPORT-TRAC":
			reader.trac = path
		case len(name) == 14 && name[:12] == "RD1002"+date[2:]:
			reader.rd1002 = path
		case len(name) > 10 && name[:9] == "IND"+date[2:] && path.UncompressedSize64 > 0:
			reader.indfiles = append(reader.indfiles, path)
		}
	}
	return
}

func (r *Reader) Check() (result bool) {
	fmt.Println("日期：", r.date)
	if r.trac == nil {
		fmt.Println("缺失 TRAC 文件")
		result = true
	}
	if r.eve == nil {
		fmt.Println("缺失 EVE 文件")
		result = true
	}
	if r.indfiles == nil {
		fmt.Println("缺失 INDFILES 文件")
		result = true
	}
	if r.rd1002 == nil {
		fmt.Println("缺失 RD1002 文件")
		result = true
	}
	return
}

func (r *Reader) Close() error {
	return r.reader.Close()
}

func loadzip(path string) (err error) {
	reader, err := OpenReader(path)
	if err != nil {
		return
	}
	defer reader.Close()

	db := Open()
	defer db.Close()

	tx, _ := db.Begin()
	defer tx.Rollback()

	reader.Check()
	reader.LoadTrac(tx)
	reader.LoadRd1002(tx)

	tx.Commit()
	return
}

func Load() {
	data_dir := filepath.Join(Home, "数据")
	pathes, _ := filepath.Glob(filepath.Join(data_dir, "??????.zip"))
	for _, path := range pathes {
		loadzip(path)
	}
}
