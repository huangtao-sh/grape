package main

import (
	"errors"
	//"grape/params/lzbg"
	"fmt"
	"grape/sqlite3"
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
	s := sqlite3.LoadSQL("insert or ignore", "lzbg", []string{"name","age"})
	fmt.Println(s)
}
