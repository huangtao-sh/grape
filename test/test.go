package main

import (
	"fmt"
	"grape/date"
)

func main() {
	d := date.Today()
	fmt.Println(d.Format("%F   %y%M%D %Y年%Q   %W %w \n%f"))
}
