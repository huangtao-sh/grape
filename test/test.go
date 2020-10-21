package main

import (
	"errors"
	"grape/data/xls"
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
	book := xls.NewFile()
	sheet := book.GetSheet("Sheet1")
	sheet.Rename("Test1")
	book.SaveAs("~/abc.xlsx")
}
