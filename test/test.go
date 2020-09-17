package main

import (
	"fmt"
	"grape/util"
)

func main() {
	fmt.Println(util.Extract(`\d{4}-\d{2}-\d{2}`, "fdsa2019-02-14sfad"))
}
