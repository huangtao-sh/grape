package main

import (
	"errors"
	//"grape/params/lzbg"
	"fmt"
	"grape/path"
	"grape/rhzh"
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
	//rhzh.LoadRhsj()
	//rhzh.LoadBhsj()
	fmt.Println("This is a test.")
	file := path.NewPath("~/Downloads").Find("单位银行结算*.xls")
	reader := rhzh.NewXlsReader(file, "PAGE1", 1)
	fmt.Println(reader)
	for rows, err := reader.Read(); err == nil; rows, err = reader.Read() {
		fmt.Println(rows)
	}

}
