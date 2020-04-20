package main

import (
	"fmt"
	"grape/sqlite"
	"grape/util"
)

func main() {
	sqlite.Config(":memory:")
	db, _ := sqlite.Open()
	defer db.Close()
	data := util.NewData()
	defer data.Wait()
	go util.PrintfCh("%s\t%02d\n", data)
	go sqlite.FetchCh(db, data, "select 'abc',12 ")
	a := sqlite.FetchValue(db, "select 'abc'")
	fmt.Println(a.(string))
}
