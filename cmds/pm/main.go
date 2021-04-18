package main

import (
	"grape"
	_ "grape/params"
	"grape/params/loadall"
	"grape/path"
	"grape/sqlite3"
)

func main() {
	path.InitLog()
	defer grape.Recover()
	defer sqlite3.Close()
	loadall.Main()
}
