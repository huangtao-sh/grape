package main

import (
	"errors"
	//"grape/params/lzbg"
	"fmt"
	"github.com/huangtao-sh/xls"
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
	if file, err := xls.Open("C:/Users/huangtao/test.xls", ""); err == nil {
		s, _ := file.GetRows("Sheet1")
		for i, r := range s {
			fmt.Println(i,len(r),r)
		}
	}
}
