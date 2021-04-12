package main

import (
	_ "grape/params"
	"grape/params/loadall"
	"grape/path"
	"grape/sqlite3"
	"grape/util"
)

func main() {
	path.InitLog()
	defer util.Recover()
	defer sqlite3.Close()
	loadall.Main()
}
