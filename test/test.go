package main

import (
	"fmt"
	"grape/date"
)

func main() {
	s := fmt.Sprint(date.Today())
	fmt.Println(s)
}
