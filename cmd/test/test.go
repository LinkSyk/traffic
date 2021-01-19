package main

import (
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	var a atomic.Value
	a.Store([]int{1})
	ch1 := make(chan struct{})
	ch2 := make(chan struct{})

	go func() {
		v := a.Load().([]int)
		v = append(v, 2)
		ch1 <- struct{}{}
		time.Sleep(5 * time.Second)
		a.Store(v)
		ch2 <- struct{}{}
	}()

	go func() {
		<-ch1
		logrus.Infoln(a.Load().([]int))
		<-ch2
		logrus.Infoln(a.Load().([]int))
	}()

}
