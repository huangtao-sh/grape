package main

import (
	"errors"
	"grape/nkwg"
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
	/*
		dt := data.NewXlsReader("C:/Users/huangtao/Downloads/resultReg.xls", 0)
		f := data.NewData()
		d := data.NewConvertReader(dt, data.Include(1,2))
		f.Add(1)
		go d.ReadAll(f)
		go f.Println()
		f.Wait()
	*/
	//nkwg.Load()
	nkwg.Report()
}
