package main

import (
	"grape/sqlite3"
)

func main() {
	sqlite3.Config("/Users/huangtao/data/abc.db")
	defer sqlite3.Close()
	sqlite3.ExecScript(`create table if not exists abc(a,b,c)`)
}
