package ggjgm

import (
	"fmt"
	"grape/data"
	"grape/params"
	"grape/sqlite3"
	"grape/text"
	"grape/util"
	"io"
)

func initDb() {
	sqlite3.ExecScript(`
-- drop table if exists ggjgm;

create table if not exists ggjgm (
    jgm text primary key,	-- 0 机构码
    mc text,		    	-- 1 机构中文名称
    jc text,				-- 3 机构简称  00、01 类型的机构存放分行的简称
    zfhh text,				-- 7 大额支付行号
    jglx text,				--15 机构类型  00-总行清算中心，01-总行营业部，10-分行清算中心，11-分行营业部，12-支行
    kbrq text,  			--16 开办日期
    hzjgm text 				--17 汇总机构码
);
	`)
}

// Reader 读取器
type Reader struct {
	*text.Reader
	io.ReadCloser
}

// NewReader 构造函数
func NewReader(file text.File) *Reader {
	r, err := file.Open()
	util.CheckFatal(err)
	reader := text.NewReader(text.Decode(r, false, true), false, text.NewSepSpliter(","),
		text.Include(0, 1, 3-43, 7-43, 15-43, 16-43, 17-43))
	return &Reader{reader, r}
}

// Exec 执行sql 语句
func (r *Reader) Exec(tx *sqlite3.Tx) (err error) {
	d := data.NewData()
	sql := "insert or replace into ggjgm values(?,?,?,?,?,date(?),?)"
	d.Add(1)
	go tx.ExecCh(sql, d)
	go r.ReadAll(d)
	d.Wait()
	return
}

// Load 导入文件
func Load(file text.File, ver string) {
	initDb()
	r := NewReader(file)
	defer r.Close()
	err := sqlite3.ExecTx(sqlite3.NewTr("delete from ggjgm"), params.NewChecker("ggjgm", file, ver), r)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("导入文件 %s 完成！\n", file.FileInfo().Name())
	}
}
