package bitmap

//用户标签工具
import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func useless23() {
	fmt.Println(1)
}

type BitMap struct {
	m            sync.RWMutex
	Data         []byte
	LastSyncTime int64 //最后同步时间
	CountUpdate  int   //最后一次同步之后又经过了多少次修改
}

func NewBitMap(data []byte) *BitMap {
	obj := &BitMap{}
	obj.Data = data
	obj.CountUpdate = 0
	obj.LastSyncTime = time.Now().Unix()
	return obj
}

func (this *BitMap) Bytes() []byte {
	return this.Data
}

//加长，如果已经够长了，什么都不操作,否则全部补充零
//注意该长度是指byte的数量长度,不是位的长度
func (this *BitMap) pad(lenth int) {
	clen := len(this.Data)
	if clen < lenth {
		d := make([]byte, lenth+lenth/5)
		copy(d, this.Data)
		this.Data = d
	}
}

//对外接口
func (this *BitMap) Get(bitpos int) bool {
	this.m.RLock()
	defer this.m.RUnlock()
	pos := bitpos / 8
	if pos > len(this.Data)-1 {
		return false
	}
	value := this.Data[pos]
	mod := bitpos % 8
	newValue := GetByte(value, mod)
	return newValue

}

//对外接口
func (this *BitMap) SetTrue(bitpos int) {
	this.Set(bitpos, true)
}
func (this *BitMap) SetFalse(bitpos int) {
	this.Set(bitpos, false)
}

func (this *BitMap) Set(bitpos int, val bool) {
	this.m.Lock()
	defer this.m.Unlock()
	pos := bitpos / 8
	if pos > len(this.Data)-1 {
		this.pad(pos + 10)
	}
	value := this.Data[pos]
	mod := bitpos % 8
	newValue := SetByte(value, mod, val)
	this.update(pos, newValue)
}
func (this *BitMap) Len() int {
	return len(this.Data) * 8
}

//update
func (this *BitMap) update(pos int, val byte) bool {
	if len(this.Data)-1 < pos {
		this.pad(pos + 10)
	}
	this.Data[pos] = val
	this.CountUpdate++
	return true

}

func minLen(t, t2 *BitMap) int {
	if len(t.Data) <= len(t2.Data) {
		return len(t.Data)
	}
	return len(t2.Data)
}

func maxLen(t, t2 *BitMap) int {
	if len(t.Data) >= len(t2.Data) {
		return len(t.Data)
	}
	return len(t2.Data)
}

func Huo(t, t2 *BitMap) *BitMap {
	len := maxLen(t, t2)
	t.pad(len)
	t2.pad(len)
	obj := &BitMap{}
	obj.Data = make([]byte, len)
	for i := 0; i < len; i++ {
		obj.Data[i] = t.Data[i] | t2.Data[i]
	}
	return obj
}

func He(t, t2 *BitMap) *BitMap {
	len := minLen(t, t2)
	obj := &BitMap{}
	obj.Data = make([]byte, len)
	for i := 0; i < len; i++ {
		obj.Data[i] = t.Data[i] & t2.Data[i]
	}
	return obj
}

func GetByte(value byte, bitpos int) bool {

	bitpos2 := byte(7 - bitpos)
	tmp := byte(1) << bitpos2
	if tmp == tmp&value {
		return true
	}

	return false

}

func SetByte(value byte, bitpos int, val bool) byte {

	bitpos2 := byte(7 - bitpos)
	if val == true {
		tmp := byte(1) << bitpos2
		if tmp == tmp&value {
			return value
		} else {
			return value | tmp
		}
	} else {
		tmp := byte(1) << bitpos2

		if byte(0) == tmp&value {
			return value
		} else {
			return value & ^tmp
		}
	}
}

//用于测试
func ByteToBinaryString(data byte) (str string) {
	var a byte
	for i := 0; i < 8; i++ {
		a = data
		data <<= 1
		data >>= 1

		switch a {
		case data:
			str += "0"
		default:
			str += "1"
		}
		data <<= 1
	}
	return str
}

func getRand() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(255)
}
