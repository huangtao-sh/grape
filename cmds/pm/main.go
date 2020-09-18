package main

import (
	_ "grape/params"
	"grape/params/loadall"
	"grape/sqlite3"
	"grape/util"
)

func main() {
	defer util.Recover()
	defer sqlite3.Close()
	loadall.Main()
}
