package main

import (
	"errors"
	"fmt"

	"github.com/extrame/xls"
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
	xl, err := xls.Open("C:/Users/huangtao/Downloads/resultReg.xls", "")
	fmt.Println(err)
	if err == nil {
		s := xl.ReadAllCells(100)
		for i, r := range s {
			fmt.Println(i, r)
		}
	}

}
