package mysqlx

// too hard to use , understand

// import(
// 	"fmt"
// 	"strconv"
// 	"os"
// )

// type Ddata struct{
// 	Data map[string]interface{}
// 	KeyName string
// 	ConStr string
// 	Sql string
// 	Done bool
// 	LastId int
// 	Step int
// 	MaxStepCount int
// }

// func NewDdata( conStr string, sql string, key string ) *Ddata{

// 	d := &Ddata{}
// 	d.Data = make(map[string]interface{} )
// 	d.KeyName = key
// 	d.ConStr = conStr
// 	d.Sql = sql
// 	d.Done = false
// 	d.LastId = 0
// 	d.Step = 1200000
// 	d.MaxStepCount = 5
// 	return d
// }

// func (this *Ddata) Loop( Lf LoopFunc ){
// 	for _ , v := range this.Data{
// 		Lf( v )
// 	}
// }

// type KeyFunc func(key interface{}) string
// type FilterFunc func( data interface{} ) bool
// type LoopFunc func( data interface{} )

// func ( this *Ddata)GroupBy( kf KeyFunc, ff FilterFunc ) map[string]int{
// 	res := make( map[string]int)
// 	fmt.Println("groupby ready")
// 	fmt.Println(len(this.Data))
// 	countGood := 0
// 	countKey := 0
// 	for _,v := range this.Data{
// 		//fmt.Println(k, v )
// 		if isGood := ff( v ); isGood{
// 	//		fmt.Println("good")
// 			key := kf( v)
// //	fmt.Println(key)
// 			if "" == key{
// 				countKey ++
// 				continue
// 			}
// 			_, ok := res[key]
// 			countGood ++
// 			if ok{
// 				res[key] ++
// 			}else{
// 				res[key] = 1
// 			}
// 		}

// 	}
// 	fmt.Println("keys",countKey)
// 	fmt.Println("goods",countGood)
// 	return res
// }

// func (this *Ddata)Merge( another *Ddata) *Ddata  {
// 	res := &Ddata{}
// 	res.Data = make(map[string]interface{} )

// 	for k,v := range this.Data{
// 		item, ok := another.Data[k]
//  		if ok{
// 			res.Data[k] = MergeMap( v, item )
// 		}
// 	}
// 	return res
// }

// func (this *Ddata )add( item map[string]interface{} , key string){
// 	this.Data[key] = item
// }

// func (this *Ddata ) Run(){

// 	for i := 1; i <= this.MaxStepCount; i++ {
// 		sql := this.Sql + " limit " + strconv.Itoa(this.LastId) + "," + " " + strconv.Itoa( this.Step)
// 		fmt.Println( sql )
// 		qr , err := Query( this.ConStr, sql, nil )
// 		if err != nil{
// 			fmt.Println( err )
// 			os.Exit(1)
// 		}
// 		this.LastId += this.Step

// 		if len( qr.M ) >0{
// 			for _, v := range qr.M{
// 				this.add( v , string( v[this.KeyName].([]uint8) ) )
// 			}
// 		}else{
// 			break
// 		}
// 	}

// }
// func (this *Ddata)Get(key string )  interface{}{
// 	item, ok := this.Data[key]
// 	if ok{
// 		return item
// 	}
// 	return nil
// }

// func (this *Ddata)Show(){
// 	fmt.Println("show",len(this.Data))

// }

// func MergeMap( m1  interface{} , m2  interface{} ) map[string]interface{}{
// 	res := make( map[string]interface{} )
// 	tmp1 := m1.(map[string]interface{})
// 	tmp2 := m2.(map[string]interface{})

// 	for k, v := range tmp1{
// 		res[k] = v
// 	}
// 	for k, v := range tmp2{
// 		res[k] = v
// 	}

// 	return res
// }
