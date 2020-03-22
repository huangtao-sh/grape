package main

import (
	"fmt"
	"grape/text"
)

func add(x interface{}) interface{} {
	return x.(int) + 10
}

func main() {
	s := []interface{}{0, 1, 2, 3, 4}
	includer := text.NewIncluder(1, 3)
	converter := text.NewConverter(map[int]text.Convert{1:add,3:add})
	fmt.Println(includer.Convert(s)...)
	fmt.Println(converter.Convert(s)...)
}
