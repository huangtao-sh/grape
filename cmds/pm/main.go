package main

import (
	"grape/params/loadall"
	"grape/util"
)

func main() {
	defer util.Recover()
	loadall.Main()
}
