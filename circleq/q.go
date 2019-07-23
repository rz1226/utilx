package circleq

import (
	"sync/atomic"
)

//只保留最近写入的部分数据

type ele struct {
	seqNo uint64
	value interface{}
}

type CQ struct {
	seqNo uint64
	size  uint32
	data  []ele
}

func NewCQ(size uint32) *CQ {
	c := &CQ{}
	c.seqNo = 0
	c.size = minQuantity(size)
	c.data = make([]ele, c.size)
	ele := &(c.data[0])
	ele.seqNo = 0
	return c
}

func (c *CQ) Put(data interface{}) {
	//func AddUint64(addr *uint64, delta uint64) (new uint64)
	newseqNo := atomic.AddUint64(&c.seqNo, 1)
	pos := newseqNo & uint64((c.size - 1))
	ele := &(c.data[pos])
	ele.seqNo = newseqNo
	ele.value = data
}

func (c *CQ) Show(len int) ([]interface{}, uint64) {
	res := make([]interface{}, len)
	seqNo := atomic.LoadUint64(&c.seqNo)
	for i := 0; i < len-1; i++ {
		pos := (seqNo - uint64(i)) & uint64((c.size - 1))
		ele := &(c.data[pos])
		if ele.seqNo == seqNo-uint64(i) {
			res[i] = ele.value
		} else {
			break
		}
	}
	return res, seqNo
}

// round 到最近的2的倍数
func minQuantity(v uint32) uint32 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}
