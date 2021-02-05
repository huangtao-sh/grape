package main

import (
	"errors"
	"fmt"
	"grape/loader"
	_ "grape/params"
	"grape/path"
	"os"
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
	file := path.NewPath("~/test.xls")
	book, err := loader.NewXlsFile(file.String())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer book.Close()
	r := book.Read(0, 0)
	var row []string
	for ; err == nil; row, err = r.Read() {
		if row != nil {
			fmt.Println(row)
		}
	}
}
