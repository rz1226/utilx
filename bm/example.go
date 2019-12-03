package bm
//
//import "fmt"
//
//type Tai struct {
//
//	Id          int64 `orm:"id" auto:"1""`
//	Name          string `orm2:"name"`
//	Age          int64 `orm2:"age"`
//	Weight float64 `orm2:"weight"`
//	Create_time string `orm2:"create_time" auto:"1"`
//	Some string
//
//}
//
//
//
//func GetOrders()  {
//	res, err := db.Kit.Query("select * from tai limit 14")
//	fmt.Println(  err )
//	//
//	u := new(Tai)
//
//	//如果字段名不一致，自己很容易调整
//	f := func(r map[string]interface{}){
//		r["name22"] = "my name is " +  fmt.Sprint(r["name"])
//	}
//
//	err2 := bm.BMFromQueryRes(res[0], u ,f  )
//	fmt.Println("err2=", err2)
//	fmt.Println("u===", u )
//
//	var uBatch  []*Tai
//	bm.BMFromQueryResBatch(res , &uBatch , f)
//	fmt.Println("batch=---------------------------------")
//	fmt.Println("ubatch=")
//	for _, v := range uBatch{
//		fmt.Println( *v )
//	}
//
//
//	fmt.Println("----------------------field line lines--------------")
//	line1 ,err := bm.LineFromBM(u)
//	fmt.Println(err)
//	fmt.Println("line1=" )
//	line1.Show()
//
//	lines ,err:= bm.LinesFromBM(uBatch)
//	fmt.Println(err )
//	fmt.Println("lines=" )
//	lines.Map("name22", func(fi *bm.Field){
//		fi.Value = "[" + fmt.Sprint(fi.Value) + "]"
//	})
//	lines.Show()
//
//	fmt.Println("---------------生成sql---------------------")
//	sql1, pa1 := bm.GetInsertSqlWithParams(line1,"tttxxx")
//	fmt.Println("sql1---")
//	fmt.Println(sql1,pa1)
//
//
//	sql, pa := bm.SqlFromLinesForInsert(lines,"tttxxx" )
//	fmt.Println("sql batch---")
//	fmt.Println(sql, pa )
//
//	fmt.Println(db.Kit.Exec(sql, pa... ))
//
//
//
//}

