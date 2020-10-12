package main

import (
	"errors"
	"fmt"
)

func getPrimes(num int) (primes []int, err error) {
	if num <= 1 {
		err = errors.New("整数不应小于1")
	} else {
		prime := 2
		for num > 1 {
			if num%prime == 0 {
				primes = append(primes, prime)
				num /= prime
			} else {
				prime++
			}
		}
	}
	return
}
func main() {
	var i int
	fmt.Printf("Please enter a number:")
	fmt.Scanf("%d", &i)
	k, err := getPrimes(i)
	if err == nil {
		fmt.Println(k)
	}
}
