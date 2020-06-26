package km

import (
	"bufio"
	"grape/gbk"
	"grape/params/load"
	"grape/text"
	"io"
	"os"
	"regexp"
	"strings"
)

var initKemuSQL = `
create table if not exists kemu(
	km      text primary key,  -- 科目
	name    text,              -- 名称
	description text           -- 说明
);
`
var loadKemuSQL = `insert or replace into kemu values(?,?,?)`

// KemuReader 科目读取
type KemuReader struct {
	*bufio.Scanner
}

// ReadAll 读取所有数据
func (r *KemuReader) ReadAll(d text.Data) {
	var line string
	KEMU := regexp.MustCompile(`^(\d{4,6})\s*(.*)`)
	BLANK1 := regexp.MustCompile(`第.章。*`)
	BLANK2 := regexp.MustCompile(`本科目为一级科目.*`)
	// AcPattern := regexp.MustCompile(`\d{1,6}`)
	defer d.Close()
	var km, name string
	var sm []string
	for r.Scan() {
		line = strings.TrimSpace(r.Text())
		if line == "" || BLANK1.MatchString(line) || BLANK2.MatchString(line) {
			continue
		} else if KEMU.MatchString(line) {
			if km != "" {
				d.Write(km, name, strings.Join(sm, "\n"))
				sm = nil
			}
			s := KEMU.FindStringSubmatch(line)
			km, name = s[1], s[2]
		} else {
			sm = append(sm, line)
		}
	}
	d.Write(km, name, strings.Join(sm, "\n"))
}

// LoadKemu 导入科目
func LoadKemu(info os.FileInfo, r io.Reader, ver string) *load.Loader {
	reader := &KemuReader{bufio.NewScanner(gbk.NewReader(r))}
	return load.NewLoader("kemu", info, ver, reader, initKemuSQL, loadKemuSQL)
}
