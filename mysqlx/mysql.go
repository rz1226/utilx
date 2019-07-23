package mysqlx

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

/*

var conf = "gechengzhen:123456@tcp(172.16.1.61:3306)/userdata?charset=utf8"
func main(){

	p := mysqlx.NewDbPool( conf,3 )

	for i := 0; i < 1000000 ; i++ {
		go test( p )
		time.Sleep( time.Millisecond * 1 )
	}


	time.Sleep( time.Second * 10000 )
}

func test( p *dbtool.DbPool ){

	 fmt.Println( p.Query("select * from ytk_car_test limit 1 ", nil ))

	//fmt.Println( p.Exec("insert into ytk_car_test set license = '赣B'", nil))
}
*/

type DbPool struct {
	pool         *sql.DB
	conStr       string
	maxOpenConns int
	maxIdleConns int
	err          error
}

func NewDbPool(conStr string, maxOpenConns int) *DbPool {
	p := &DbPool{}
	p.conStr = conStr
	p.pool, p.err = sql.Open("mysql", conStr)
	if p.err == nil {
		p.pool.SetMaxOpenConns(maxOpenConns)
		p.pool.SetMaxIdleConns(maxOpenConns)
		p.pool.SetConnMaxLifetime(time.Second * 10000)
	} else {
		p.pool = nil
	}
	return p
}

//获取*sql.DB
func (p *DbPool) DB() *sql.DB {
	return p.pool
}

func (p *DbPool) Query(sqlStr string, data []interface{}) (queryResult, error) {
	rows, err := p.query(sqlStr, data)
	if err == nil {
		return rowsToMap(rows)
	} else {
		return queryResult{nil, nil}, err
	}

}

func (p *DbPool) query(sqlStr string, data []interface{}) (*sql.Rows, error) {
	if p.pool == nil {
		return nil, p.err
	}
	db := p.pool
	len := len(data)
	fn := reflect.ValueOf(db.Query)
	params := make([]reflect.Value, len+1)
	params[0] = reflect.ValueOf(sqlStr)
	for i := 1; i <= len; i++ {
		params[i] = reflect.ValueOf(data[i-1])
	}
	//fmt.Println( params )
	fv := fn.Call(params)
	if fv[1].Interface() != nil {
		return nil, fv[1].Interface().(error)
	}
	return fv[0].Interface().(*sql.Rows), nil
}

// query返回的结果
type queryResult struct {
	M []map[string]interface{}
	S [][]interface{}
}

//把数据库取出来d数据的rows的数据放在一个queryResult上
// null 对应nil  数字对数字  其他对字符串

func rowsToMap(rows *sql.Rows) (queryResult, error) {
	defer rows.Close()
	res := make([]map[string]interface{}, 0)
	s := make([][]interface{}, 0)
	fields, err := rows.Columns()
	len := len(fields)
	if err != nil {
		return queryResult{nil, nil}, err
	}

	for {
		if result := rows.Next(); result {
			data := make([]interface{}, len)
			data2 := make(map[string]interface{}, len)
			dvalue := reflect.ValueOf(&data)
			fn := reflect.ValueOf(rows.Scan)
			params := make([]reflect.Value, len)
			for i := 0; i < len; i++ {
				params[i] = dvalue.Elem().Index(i).Addr()
			}
			fv := fn.Call(params)
			if fv[0].Interface() != nil {
				return queryResult{nil, nil}, fv[0].Interface().(error)
			}
			for i := 0; i < len; i++ {
				data2[fields[i]] = data[i]
			}
			res = append(res, data2)
			s = append(s, data)
		} else {
			break
		}
	}
	qr := queryResult{}
	qr.M = res
	qr.S = s
	return qr, nil
}

//datas是个二维slice,
//注意，这仍然是每个条数据遍历执行，插入的时候并不是自动生成批量sql插入的。
func (p *DbPool) ExecMany(sqlStr string, datas [][]interface{}) (int64, error) {
	if p.pool == nil {
		return 0, p.err
	}
	db := p.pool
	var rowsAffected int64 = 0
	var lastInsertId int64 = 0

	//插入数据
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	if err != nil {
		return 0, err
	}
	fn := reflect.ValueOf(stmt.Exec)
	for _, data := range datas {
		len := len(data)
		params := make([]reflect.Value, len)
		for i := 0; i < len; i++ {
			params[i] = reflect.ValueOf(data[i])
		}
		fv := fn.Call(params)
		if fv[1].Interface() != nil {
			continue
			//return 0, fv[1].Interface().(error)
		}
		result := fv[0].Interface().(sql.Result)
		if isUpdate(sqlStr) || isDelete(sqlStr) {
			tmpAffetct, err := result.RowsAffected()
			if err != nil {
				continue
			}
			rowsAffected += tmpAffetct
		}
		if isInsert(sqlStr) {
			tmpLast, err := result.LastInsertId()
			if err != nil {
				continue
			}
			lastInsertId = tmpLast
		}
	}
	if isUpdate(sqlStr) || isDelete(sqlStr) {
		return rowsAffected, nil
	}
	if isInsert(sqlStr) {
		return lastInsertId, nil
	}
	return 0, errors.New("only support update insert delete ")
}

// data 是一个slice, 里面的个数对应sqlStr里面？的数量
func (p *DbPool) Exec(sqlStr string, data []interface{}) (int64, error) {
	if p.pool == nil {
		return 0, p.err
	}
	db := p.pool
	//插入数据
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	len := len(data)
	if err != nil {
		return 0, err
	}

	fn := reflect.ValueOf(stmt.Exec)
	params := make([]reflect.Value, len)
	for i := 0; i < len; i++ {
		params[i] = reflect.ValueOf(data[i])
	}
	fv := fn.Call(params)

	if fv[1].Interface() != nil {
		return 0, fv[1].Interface().(error)
	}
	result := fv[0].Interface().(sql.Result)
	if isUpdate(sqlStr) || isDelete(sqlStr) {
		return result.RowsAffected() //本身就是多个返回值
	}
	if isInsert(sqlStr) {
		return result.LastInsertId() //本身就是多个返回值
	}
	return 0, errors.New("only support update insert delete ")
}

func isInsert(sqlStr string) bool {
	str := strings.TrimSpace(strings.ToLower(sqlStr))
	if strings.HasPrefix(str, "insert") {
		return true
	}
	return false
}

func isUpdate(sqlStr string) bool {
	str := strings.TrimSpace(strings.ToLower(sqlStr))
	if strings.HasPrefix(str, "update") {
		return true
	}
	return false
}

func isDelete(sqlStr string) bool {
	str := strings.TrimSpace(strings.ToLower(sqlStr))
	if strings.HasPrefix(str, "delete") {
		return true
	}
	return false
}

//
func Int64(data interface{}) int64 {
	if data == nil {
		return 0
	}
	str := string(data.([]uint8))
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		fmt.Println("int64 conv error")
		log.Fatal(err)
	}
	return num
}

func Int(data interface{}) int64 {
	if data == nil {
		return 0
	}
	str := string(data.([]uint8))
	num, err := strconv.ParseInt(str, 10, 0)
	if err != nil {
		fmt.Println("int conv error")
		log.Fatal(err)
	}
	return num
}

func String(data interface{}) string {
	if data == nil {
		return ""
	}
	str := string(data.([]uint8))
	return str
}

func Float64(data interface{}) float64 {
	if data == nil {
		return 0
	}
	str := string(data.([]uint8))
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		fmt.Println(err)
	}
	return f

}
