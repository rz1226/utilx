package actor

import (
	"fmt"
	"github.com/rz1226/simplegokit/blackboardkit"
	"github.com/rz1226/simplegokit/coroutinekit"
	"sync/atomic"
)


/**
a := actor.NewActor(inf ,1 )  第一个actor初始化
b := a.AddActor(inf2 ,2)  后续的actor初始化 并且连起来

a.Run()  启动
for{
	a.Put(xxx) 持续投入参数
}
*/
var currentId uint64 = 0
var actorbb *blackboardkit.BlackBoradKit

func init() {
	actorbb = blackboardkit.NewBlockBorad()
	actorbb.InitLogKit("actor_error")
	actorbb.SetLogReadme("actor_error", "最近错误日志")
	actorbb.InitCounterKit("actor")
	actorbb.SetCounterReadme("actor", "actor执行总次数")
	actorbb.InitTimerKit("actor_function_run", )
	actorbb.SetTimerReadme("actor_function_run", "actor耗时计数")
	actorbb.SetNoPrintToConsole(true)
	actorbb.SetName("-------actor------")
	actorbb.Ready()
}

//模拟的actor
type Actor struct{
	Id uint64
	C chan interface{}
	F func(interface{})(interface{} ,error)
	NumOfConcurrent uint8  //并发数量
	Next []*Actor
}

func NewActor( f func(interface{})(interface{} , error ), num int ) ( *Actor){
	a := &Actor{}
	a.Id = atomic.AddUint64(&currentId,1)
	a.C = make( chan interface{} ,0 )
	a.F = f
	a.NumOfConcurrent = uint8(num)
	if a.NumOfConcurrent <= 0 {
		a.NumOfConcurrent = 1
	}
	a.Next = make([]*Actor , 0 ,10)
	return a
}
func ( a *Actor) AddActor(f func(interface{})(interface{},error) ,num int )*Actor {
	b := NewActor(f , num  )
	a.setNext( b )
	return b
}

func (a *Actor) setNext( b *Actor ){
	a.Next = append( a.Next, b )
}

func (a *Actor)Run(){
	a.run()
}
func (a *Actor)Put(data interface{}) {
	a.C <- data
}

func (a *Actor)run(){
	workF := func(){
		for{
			//从队列获取数据
			data :=  <- a.C
			if data == nil{
				continue
			}
			if a.F != nil {
				actorbb.Inc("actor")
				t := actorbb.Start("actor_function_run","actor_function_run")
				res , err := a.F( data )
				actorbb.End(t )
				if err != nil {
					actorbb.Log("actor_error", err)
					continue
				}
				for _ , v := range a.Next{
					v.Put(res )
				}
			}else{
				//如果f没有设置，或者为nil,直接不处理输出到出口chan
				for _ , v := range a.Next{
					v.Put(data )
				}
			}
		}
	}

	coroutinekit.Start("actor job"+ fmt.Sprint(a.Id ), int(a.NumOfConcurrent), workF, true)
	if len(a.Next) > 0    {
		for _, v := range a.Next{
			v.run()
		}
	}
}


