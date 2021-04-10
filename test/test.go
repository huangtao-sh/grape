package main

import (
	"errors"
	"fmt"
	_ "grape/params"
	"grape/path"
	"grape/util"
	//"grape/params/lzbg"
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
	file := path.NewPath("~/Documents/参数备份/营业主管").Find("营业主管信息*.xls*")
	info := path.NewPath(file).FileInfo()
	s := util.Extract("\\d+", info.Name())
	fmt.Println(s,info.Name())
}
