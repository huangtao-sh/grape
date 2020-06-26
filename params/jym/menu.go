package jym

import (
	"encoding/xml"
	"errors"
	"grape/gbk"
	"grape/params/load"
	"grape/text"
	"io"
	"os"
)

func read(charset string, r io.Reader) (result io.Reader, err error) {
	switch charset {
	case "GBK":
		result = gbk.NewReader(r)
	case "UTF8":
		result = r
	default:
		err = errors.New("无法解码")
	}
	return
}

func getAttr(token xml.StartElement) (attrs map[string]string) {
	attrs = make(map[string]string)
	for _, a := range token.Attr {
		attrs[a.Name.Local] = a.Value
	}
	return
}

var initSQL = `
create table if not exists menu(
	jym		text	primary key, -- 交易码
	name	text,	-- 交易名称
	yjcd	text, 	-- 一级菜单
	ejcd	text	-- 二级菜单
)`
var loadSQL = `insert or replace into menu values(?,?,?,?)`

// LoadMenu 导入科目
func LoadMenu(info os.FileInfo, r io.Reader, ver string) *load.Loader {
	reader := &MenuReader{r}
	return load.NewLoader("menu", info, ver, reader, initSQL, loadSQL)
}

// MenuReader 读取菜单文件
type MenuReader struct {
	r io.Reader
}

// ReadAll 读取数据
func (l *MenuReader) ReadAll(dt text.Data) {
	defer dt.Close()
	d := xml.NewDecoder(l.r)
	d.CharsetReader = read
	submenu := make([]string, 2, 2)
	i := -1
	for t, err := d.Token(); err == nil; t, err = d.Token() {
		switch token := t.(type) {
		case xml.StartElement:
			switch token.Name.Local {
			case "SubMenu":
				attr := getAttr(token)
				i++
				submenu[i] = attr["DisplayName"]
			case "Trade":
				attr := getAttr(token)
				dt.Write(attr["Code"], attr["DisplayName"], submenu[0], submenu[1])
			}
		case xml.EndElement:
			if token.Name.Local == "SubMenu" {
				submenu[i] = ""
				i--
			}
		}
	}
}

/*
// MenuLoader  菜单文件导入
type MenuLoader struct {
	file string
	data *data.Data
}

// NewMenuLoader 构造 MenuLoader
func NewMenuLoader(file string) *MenuLoader {
	return &MenuLoader{file, data.NewData()}
}

// Read 读取数据
func (l *MenuLoader) Read() {
	f, err := os.Open(l.file)
	util.CheckFatal(err)
	defer f.Close()
	d := xml.NewDecoder(f)
	d.CharsetReader = read
	submenu := []string{}
	defer l.data.Close()
	for t, err := d.Token(); err == nil; t, err = d.Token() {
		switch token := t.(type) {
		case xml.StartElement:
			switch token.Name.Local {
			case "SubMenu":
				attr := getAttr(token)
				submenu = append(submenu, attr["DisplayName"])
			case "Trade":
				attr := getAttr(token)
				l.data.Write(attr["Code"], attr["DisplayName"], strings.Join(submenu, "/"))
			}
		case xml.EndElement:
			if token.Name.Local == "SubMenu" {
				submenu = submenu[0 : len(submenu)-1]
			}
		}
	}
}

// Exec 执行程序
func (l *MenuLoader) Exec(tx *sqlite3.Tx) (err error) {
	tx.Exec("delete from menu")
	l.data.Add(1)
	go l.data.Exec(tx, "insert or replace into menu values(?,?,?)")
	go l.Read()
	l.data.Wait()
	fmt.Println("导入完成")
	return
}

// Load 导入数据
func (l *MenuLoader) Load() {
	initDb()
	sqlite3.ExecTx(l)
}

// Test 测试数据
func (l *MenuLoader) Test() {
	l.data.Add(1)
	go l.data.Println()
	go l.Read()
	l.data.Wait()
	fmt.Println("测试完成")
}

*/
