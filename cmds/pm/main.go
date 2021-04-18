package main

import (
	_ "grape/params"
	"grape/params/loadall"
	"grape"
	"grape/sqlite3"
)

func main() {
	grape.InitLog()
	defer grape.Recover()
	defer sqlite3.Close()
	loadall.Main()
}
