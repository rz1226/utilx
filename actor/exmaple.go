package actor

//
//import (
//	"errors"
//	"fmt"
//	"github.com/rz1226/simplegokit/coroutinekit"
//	"sync/atomic"
//
//	"github.com/rz1226/simplegokit/blackboardkit"
//	"time"
//)
//func init() {
//
//	coroutinekit.StartMonitor("9090")
//	blackboardkit.StartMonitor("9091")
//
//}
//func main(){
//	f := func(data interface{}) (interface{},error){
//		s := fmt.Sprint(data ) + "_f"
//		return s, nil
//	}
//	f2 := func(data interface{}) (interface{},error){
//		s := fmt.Sprint(data ) + "_f2"
//		return s, nil
//	}
//	f3 := func(data interface{}) (interface{},error){
//		fmt.Println(data )
//		return nil, nil
//	}
//	var  count uint64 = 0
//	f4 := func(data interface{}) (interface{},error){
//		atomic.AddUint64(&count,1)
//		if atomic.LoadUint64(&count) % 10 == 0 {
//			return "累计" + fmt.Sprint(atomic.LoadUint64(&count)), nil
//		}else{
//			//如果返回nil，不再往下一层传递
//			return nil, errors.New("示例错误")
//		}
//
//	}
//	a :=  NewActor(nil ,2  ) ;
//	c := a.AddActor(f,2).AddActor(f2 ,3) ;
//	c.AddActor(f3,2)
//	c.AddActor(f4,3).AddActor(f3,1)
//
//	a.Run() ;
//	for i:=0;i<2000070;i++{
//
//		a.Put(i)
//	}
//
//
//	time.Sleep(time.Second*100)
//}
//
