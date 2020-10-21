package main

import (
	"errors"
	"fmt"
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
	var widthes = []struct {
		Name  string
		Width float64
	}{
		{"A", 15},
		{"B", 45},
	}
	for _, v := range widthes {
		fmt.Println(v.Name, v.Width)
	}
}
