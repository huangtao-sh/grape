package main

import (
	"errors"
	"fmt"
	"grape/data/xls"
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

// Test a test class
type Test struct {
	Name string `json:"姓名"`
	Age  int64  `json:"年龄"`
}

func main() {
	sqlite3.Config("params")
	defer sqlite3.Close()
	/*
		widthes := map[string]float64{
			"A:D": 40,
		}
		sql := "select ygh,xm,js,lxdh from yyzg limit 10"
	*/
	book := xls.NewFile()
	style1, err := book.NewStyle(`{"border":[]}`)
	if err != nil {
		fmt.Println(err)
	}
	book.SaveAs("~/test.xlsx")

}
