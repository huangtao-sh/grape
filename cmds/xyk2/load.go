package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"database/sql"
	"fmt"
	"grape"
	"grape/sqlite3"
	"grape/text"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

func GetDate() (rq string) {
	sqlite3.QueryRow("select max(rq) from ylqs").Scan(&rq)
	return
}

func findPath(p string, pattern string) *grape.Path {
	a := grape.NewPath(p).Glob(pattern)
	if len(a) >= 1 {
		return grape.NewPath(a[0])
	}
	fmt.Printf("can not find file %s ", pattern)
	os.Exit(1)
	return nil
}

func Load() {
	root := grape.NewPath("~/信用卡")
	tx, err := sqlite3.NewDB().Begin()
	grape.CheckFatal(err)
	defer tx.Rollback()
	date := GetDate()
	wg := &sync.WaitGroup{}
	for _, p := range root.Glob("??????") {
		rq := grape.NewPath(p).Base()
		if rq > date {
			fmt.Println("处理文件夹 ", p)
			rd1002 := findPath(p, "RD1002??????99")
			trac := findPath(p, "GLREPORT-TRAC*")
			wg.Add(2)
			go loadRd1002(tx, rq, rd1002, wg)
			go loadTrac(tx, rq, trac, wg)
		}
	}
	wg.Wait()
	tx.Commit()
}

func NewScanner(p *grape.Path) *bufio.Scanner {
	var r io.Reader
	var err error
	b, err := ioutil.ReadFile(p.String())
	grape.CheckFatal(err)
	r = bytes.NewReader(b)
	if p.Ext() == ".gz" {
		r, err = gzip.NewReader(r)
		grape.CheckFatal(err)
	}
	r = grape.NewGBKReader(r)
	return bufio.NewScanner(r)
}

func loadTrac(tx *sql.Tx, rq string, p *grape.Path, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := NewScanner(p)
	for scanner.Scan() {
		s := scanner.Text()
		if strings.Contains(s, "440001") {
			fields := strings.Fields(s)
			_, err := tx.Exec("insert or replace into ys values(?,?,?)", rq, fields[3], fields[4])
			grape.CheckFatal(err)
		}
	}
}

func loadRd1002(tx *sql.Tx, rq string, p *grape.Path, wg *sync.WaitGroup) {
	defer wg.Done()
	var fields []string
	var err error
	scanner := NewScanner(p)
	sff := []string{"0", "0"}
	d := []string{rq}
	for scanner.Scan() {
		s := scanner.Text()
		if strings.Contains(s, "收付费") {
			fields := strings.Fields(s)
			sff = fields[2:4]
		} else if strings.Contains(s, "总    计") {
			fields = strings.Fields(s)
			d = append(d, sff...)
			d = append(d, fields[2:]...)
		}
	}
	data := text.Slice(d)
	_, err = tx.Exec(`insert or replace into ylqs Values(?,?,?,?,?,?,?,?,?,?,?,?)`, data...)
	grape.CheckFatal(err)
}
