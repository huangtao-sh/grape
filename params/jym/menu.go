package jym

import (
	"encoding/xml"
	"errors"
	"grape/params/load"
	"grape/text"
	"io"
	"grape"
	"os"
)

func read(charset string, r io.Reader) (result io.Reader, err error) {
	switch charset {
	case "GBK":
		result = grape.NewGBKReader(r)
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

const loadSQL = `insert or replace into menu values(?,?,?,?)`

// LoadMenu 导入科目
func LoadMenu(info os.FileInfo, r io.Reader, ver string) *load.Loader {
	InitJym()
	reader := &MenuReader{r}
	return load.NewLoader("menu", info, ver, reader, "", loadSQL)
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
	submenu := make([]string, 2)
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
