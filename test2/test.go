package main

import (
	"fmt"
	"sync"
)

type DataCh struct {
	sync.WaitGroup
	ch chan []interface{}
}

func NewDataCh() *DataCh {
	ch := make(chan []interface{})
	return &DataCh{ch: ch}
}

func (d *DataCh) ReadCh() <-chan []interface{} {
	return d.ch
}

func (d *DataCh) WriteCh() chan<- []interface{} {
	return d.ch
}

func (d *DataCh) Close() {
	close(d.ch)
}

func Print(d *DataCh) {
	defer d.Done()
	for i := range d.ReadCh() {
		fmt.Println(i)
	}
}
func Create(d *DataCh) {
	defer d.Close()
	for i := 0; i < 100; i++ {
		d.WriteCh() <- []interface{}{i + 1, i*10 + 1}
	}
}

func main() {
	d := NewDataCh()
	d.Add(2)
	go func(d *DataCh) {
		defer d.Done()
		for i := range d.ReadCh() {
			fmt.Println("hello", i)
		}
	}(d)
	go func(d *DataCh) {
		defer d.Done()
		for i := range d.ReadCh() {
			fmt.Println("world", i)
		}
	}(d)
	go Create(d)
	d.Wait()
}
