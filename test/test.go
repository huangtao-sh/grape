package main

import (
	"fmt"
	"grape/path"
	"grape/sqlite"
)

func init() {
	sqlite.Config("test.db")
}

func main() {
	db, err := sqlite.Open()
	if err != nil {
		fmt.Println("Fatal")
		return
	}
	defer db.Close()
	db.Exec("create table if not exists abc(a,b)")
	fmt.Println("Hellow")
	abc := path.NewPath("~")
	d, _ := abc.Glob("/Music/*")
	for _,k :=range d{
		fmt.Println(k)
	}
}
