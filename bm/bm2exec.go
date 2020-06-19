package bm

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// 如果想要灵活的生成insert sql, 不应该用struct去生成，因为这里面用到了struct的tag，field，这些是无法通过直接代码改变的。
// 从struct改为一种容易直接修改的数据结构Line , Lines来生成sql。
//{field = name, auto : 1 , value = xx}, ...
// 一个拍扁了的orm struct
type Field struct {
	Name   string      // 数据库字段名
	IsAuto bool        //是否类似auto_increment, 或者create_time 不需要手动设置的数据
	Value  interface{} //值
}
type Lines []Line

func (l Lines) Map(fieldName string, f func(*Field)) {
	for _, v := range l {
		field, err := v.GetField(fieldName)
		if err == nil {
			f(field)
		}
	}
}

// key 是字段名
type Line []*Field

func (e Line) GetField(field string) (*Field, error) {
	for _, v := range e {
		if v.Name == field {
			return v, nil
		}
	}
	return nil, errors.New("not found field :" + field)
}
func (e Line) Show() {
	for k, v := range e {
		fmt.Println("field index=", k, " v=", *v)
	}
}
func (e Lines) Show() {
	for k, v := range e {
		fmt.Println("line index:", k)
		v.Show()
	}
}

/**
这个结构体叫做 业务模型 BM
type Tai struct {
	Id          int64 `orm:"id" auto:"1""`
	Name          string `orm:"name22"`
	Age          int64 `orm:"age"`
	Weight float64 `orm:"weight"`
	Create_time string `orm:"create_time" auto:"1"`
}

*/
// 参数类型直接是struct or 其指针
func LineFromBM(sourceStruct interface{}) (Line, error) {
	v := reflect.ValueOf(sourceStruct)
	res := make([]*Field, 0, 20)
	t := v.Type()
	switch v.Kind() {
	case reflect.Ptr:
		return LineFromBM(v.Elem().Interface())
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fd := &Field{}
			fieldName := t.Field(i).Tag.Get(Conf.TagName)
			if fieldName == "" {
				//没找到映射tag 直接忽略
				continue
			}
			value := v.Field(i).Interface()
			fd.Value = value
			fd.Name = fieldName
			fd.IsAuto = Conf.FuncAuto(t.Field(i).Tag)
			res = append(res, fd)
		}
	default:
		return res, errors.New("LineFromBM 参数1 必须是struct或者其指针")
	}
	return res, nil
}

//上面函数的批量方法
// 参数是[]*struct  or []struct
func LinesFromBM(sourceStructArray interface{}) (Lines, error) {
	res := make([]Line, 0, 100)
	v := reflect.ValueOf(sourceStructArray)
	length := v.Len()

	switch v.Kind() {
	case reflect.Ptr:
		return LinesFromBM(v.Elem().Interface())
	case reflect.Slice:
		for i := 0; i < length; i++ {
			ele := v.Index(i).Interface()
			line, err := LineFromBM(ele)
			if err != nil {
				return res, err
			}
			res = append(res, line)
		}
	}
	return res, nil
}

//auto代表插入的时候不管  影响生成insert语句
//type User struct{
//	Id int64 `orm:"id" auto:"1"`
//	Phone string `orm:"phone"`
//	Passwd string `orm:"passwd"
//	CreateTime string `orm:"create_time" auto:"1"`
//}

// 1 insertFieldList , 2 insert ? 占位符 3 insert parmas
// 如果想要灵活的生成insert sql, 不应该用struct去生成，因为这里面用到了struct的tag，field，这些是无法通过直接代码改变的。
// 所以这里要变，把第一个参数从struct改为一种容易直接修改的数据结构。
func sqlFromLine(d Line) (string, string, []interface{}) {
	insertFieldList := "("
	insertMarksStr := "("
	insertValuesSli := make([]interface{}, 0, 30)
	//type lineData map[string]*fieldData
	for _, v := range d {
		//略过id

		auto := v.IsAuto
		if auto == false {
			insertFieldList += "`" + v.Name + "`" + ","
			insertMarksStr += "?,"
			insertValuesSli = append(insertValuesSli, v.Value)
		}
	}
	return strings.TrimRight(insertFieldList, ",") + ")",
		strings.TrimRight(insertMarksStr, ",") + ")",
		insertValuesSli
}

//生成一个insert语句
func GetInsertSqlWithParams(d Line, tableName string) (string, []interface{}) {
	insertFields, insertMarks, insertParams := sqlFromLine(d)
	insertSql := "insert into " + tableName + "  " + insertFields + " values " + insertMarks
	return insertSql, insertParams
}
func SqlFromLinesForInsert(d Lines, tableName string) (string, []interface{}) {
	length := len(d)
	if length == 0 {
		return "", nil
	}
	var marksBuf bytes.Buffer
	insertFields, insertMarks, insertParams := sqlFromLine(d[0])
	marksBuf.WriteString(insertMarks)
	marksBuf.WriteString(",")
	for i := 1; i < length; i++ {
		_, marks, params := sqlFromLine(d[i])
		marksBuf.WriteString(marks)
		marksBuf.WriteString(",")
		insertParams = append(insertParams, params...)
	}
	insertSql := "insert into " + tableName + "  " + insertFields + " values " + strings.Trim(marksBuf.String(), ",")
	return insertSql, insertParams
}
