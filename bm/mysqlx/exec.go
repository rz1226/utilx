package mysqlx

import (
	"database/sql"
	"errors"
	"strings"
)

// exec 用来执行非query sql
func (p *DbPool) Exec(sqlStr string, args ...interface{}) (int64, error) {
	return execCommon(p, sqlStr, args...)
}
func (t *DbTx) Exec(sqlStr string, args ...interface{}) (int64, error) {
	return execCommon(t, sqlStr, args...)
}

func execCommon(source interface{}, sqlStr string, args ...interface{}) (int64, error) {
	p, ok := source.(*DbPool)
	if ok {
		result, err := p.realPool.Exec(sqlStr, args...)
		if err != nil {
			return int64(0), err
		}
		return affectedResult(sqlStr, result)
	}
	t, ok := source.(*DbTx)
	if ok {
		result, err := t.realtx.Exec(sqlStr, args...)
		if err != nil {
			return int64(0), err
		}
		return affectedResult(sqlStr, result)
	}
	return int64(0), errors.New("only support DbPool , DbTx")
}

//从exec的result获取   当insert获取最后一个id， update，delete获取影响行数，replace获取最后一个id
func affectedResult(sqlStr string, result sql.Result) (int64, error) {
	if isSqlUpdate(sqlStr) || isSqlDelete(sqlStr) {
		return result.RowsAffected() //本身就是多个返回值
	}
	if isSqlInsert(sqlStr) {
		return result.LastInsertId() //本身就是多个返回值
	}
	if isSqlReplace(sqlStr) {
		return result.LastInsertId() //本身就是多个返回值
	}
	return int64(0), errors.New("only support update insert delete replace")
}

func isSqlReplace(sqlStr string) bool {
	str := strings.TrimSpace(strings.ToLower(sqlStr))
	if strings.HasPrefix(str, "replace") {
		return true
	}
	return false
}
func isSqlInsert(sqlStr string) bool {
	str := strings.TrimSpace(strings.ToLower(sqlStr))
	if strings.HasPrefix(str, "insert") {
		return true
	}
	return false
}

func isSqlUpdate(sqlStr string) bool {
	str := strings.TrimSpace(strings.ToLower(sqlStr))
	if strings.HasPrefix(str, "update") {
		return true
	}
	return false
}

func isSqlDelete(sqlStr string) bool {
	str := strings.TrimSpace(strings.ToLower(sqlStr))
	if strings.HasPrefix(str, "delete") {
		return true
	}
	return false
}