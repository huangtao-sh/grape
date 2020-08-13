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

// Home 本程序的家目录
var Home string

// 初始化根目录
func init() {
	home, _ := os.UserHomeDir()
	Home = filepath.Join(home, "信用卡对账")
}

// Reader  zip 文件读取
type Reader struct {
	reader   *zip.ReadCloser
	date     string      // 日期
	trac     *zip.File   // 银数日报表
	jorj     *zip.File   // 银数会计流水
	eve      *zip.File   // 银数交易日志
	rd1002   *zip.File   // 银联日报表
	indfiles []*zip.File // 银联明细文件
}

// ReadAll 读取压缩包内文件
func ReadAll(file *zip.File, isGbk bool) (data []byte, err error) {
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
	if isGbk {
		reader = gbk.NewReader(reader)
	}
	return ioutil.ReadAll(reader)
}

// OpenReader 打开压缩包
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

// CheckFileList 检查文件列表是否齐全
func (r *Reader) CheckFileList() (result bool) {
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

// Close 关闭压缩包
func (r *Reader) Close() error {
	return r.reader.Close()
}

// loadzip 打开 zip 文件
func loadzip(path string) (err error) {
	reader, err := OpenReader(path)
	if err != nil {
		return
	}
	defer reader.Close()
	if reader.CheckFileList() { // 检查文件列表
		panic("文件不全")
	}
	db := Open() // 打开数据库连接
	defer db.Close()
	tx, _ := db.Begin() // 开启事务
	defer tx.Rollback()

	reader.LoadTrac(tx)
	reader.LoadRd1002(tx)
	reader.LoadJorj(tx)
	reader.LoadEve(tx)
	reader.LoadInds(tx)
	reader.ChongZheng(tx)
	tx.Commit()
	return
}

// Load 导入文件
func Load() {
	dataDir := filepath.Join(Home, "数据")
	pathes, _ := filepath.Glob(filepath.Join(dataDir, "??????.zip"))
	for _, path := range pathes {
		loadzip(path)
	}
}
