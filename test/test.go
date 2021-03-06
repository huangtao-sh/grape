package main

import (
	"errors"
	//"grape/params/lzbg"
	"fmt"
	"grape"
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
	s := grape.ExtractPos(`(结算账户|电子结算户)\((.*?)\)`,"结算账户()", 2)
	fmt.Println(s)

}
