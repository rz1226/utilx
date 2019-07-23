package numlimit

import (
	"fmt"
	"log"
	"sync/atomic"
)

/*
var x := NewNumLimiter( 5 )
x.Add()  //超过限制返回false ,否则返回true
x.Done()

*/

type NumLimiter struct {
	max   uint32
	value uint32
}

func NewNumLimiter(max uint32) *NumLimiter {
	n := &NumLimiter{}
	n.max = max
	n.value = 0
	return n
}

func (n *NumLimiter) Done() {
	if n.value == 0 {
		log.Println("numlimit ,this should not happen")
		return
	}
	atomic.AddUint32(&n.value, ^uint32(0))
}

//并发的问题可以容忍
func (n *NumLimiter) Add() bool {

	if n.value >= n.max {
		return false
	} else {
		atomic.AddUint32(&n.value, 1)
		return true
	}

}
func (n *NumLimiter) Show() {
	fmt.Println(n.value)
}
