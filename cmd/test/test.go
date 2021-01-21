package main

import (
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	var a atomic.Value
	a.Store(map[int]string{
		1: "a",
		2: "b",
	})
	ch1 := make(chan struct{})
	ch2 := make(chan struct{})

	go func() {
		v := a.Load().(map[int]string)
		v[3] = "c"
		ch1 <- struct{}{}
		time.Sleep(5 * time.Second)
		a.Store(v)
		ch2 <- struct{}{}
	}()

	go func() {
		<-ch1
		logrus.Infoln(a.Load().(map[int]string))
		<-ch2
		logrus.Infoln(a.Load().(map[int]string))
	}()

	for {
		time.Sleep(time.Second * 1)
	}
}
