package main

import (
	"errors"
	"fmt"
	"grape/loader"
	_ "grape/params"
	"grape/path"
)

func getPrimes(num int) (primes []int, err error) {
	if num <= 1 {
		err = errors.New("整数不应小于1")
		return
	}
	prime := 2
	for num > 1 {
		if num%prime == 0 {
			primes = append(primes, prime)
			num /= prime
		} else {
			prime++
		}
	}
	return
}

func main() {
	file := path.NewPath("~/Downloads/resultReg (1).xls")
	r := loader.NewXlsReader(file.String(), 0, 1)
	r = loader.NewConverter(r, loader.Include(0, 1))
	var row []string
	var err error
	var i int
	for ; err == nil && i<2; row, err = r.Read() {
		if row != nil  {
			fmt.Println(loader.Slice(row)...)
			i++
		}
	}

}
