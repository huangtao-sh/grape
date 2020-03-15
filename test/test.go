package main

import (
	"bytes"
	"grape/text"
	"grape/util"
)

func main() {
	s := `abcddef
testsfd
hwerwer`
	b := bytes.NewBufferString(s)
	csv := text.NewFixedReader(b, []int{0, 4, 7})
	util.Println(csv)
}
