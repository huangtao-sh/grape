package trans

import (
	"encoding/xml"
	"errors"
	"fmt"
	"grape/data"
	"grape/gbk"
	"grape/sqlite3"
	"grape/util"
	"io"
	"os"
	"strings"
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

func initDb() {
	sqlite3.ExecScript(`
create table if not exists menu(
	jym		text	primary key,
	name	text,
	menu	text
)
	`)
	fmt.Println("建立表结构完成")
}

// MenuLoader  菜单文件导入
type MenuLoader struct {
	file string
	data *data.Data
}

// NewMenuLoader 构造 MenuLoader
func NewMenuLoader(file string) *MenuLoader {
	return &MenuLoader{file, data.NewData()}
}

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
func (l *MenuLoader) Exec(tx *sqlite3.Tx) (err error) {
	tx.Exec("delete from menu")
	l.data.Add(1)
	go l.data.Exec(tx, "insert or replace into menu values(?,?,?)")
	go l.Read()
	l.data.Wait()
	fmt.Println("导入完成")
	return
}
func (l *MenuLoader) Load() {
	initDb()
	sqlite3.ExecTx(l)
}

func (l *MenuLoader) Test() {
	l.data.Add(1)
	go l.data.Println()
	go l.Read()
	l.data.Wait()
	fmt.Println("测试完成")
}
