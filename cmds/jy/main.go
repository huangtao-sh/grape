package main

import (
	"grape/params/jym"
	"grape/sqlite3"
	"grape/util"
)

func main() {
	defer util.Recover()
	defer sqlite3.Close()
	jym.Main()
}
