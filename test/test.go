package main

import (
	"errors"
	"fmt"
)

func yz(num int) (result []int, err error) {
	if num <= 1 {
		err = errors.New("整数不应小于1")
	} else {
		prime := 2
		for num > 1 {
			if num%prime == 0 {
				result = append(result, prime)
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
	k, err := yz(i)
	if err == nil {
		fmt.Println(k)
	}
}
