package main

import (
	"grape/params/wwxt"
	"grape/util"
)

func main() {
	defer util.Recover()
	wwxt.Main()
}
