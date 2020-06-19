package mysqlx

import (
	"database/sql"
	"time"
)

/**
基本思想是这样的，暴露每一步和下一步之间的数据结构给用户，让用户可以灵活修改。
例如有的框架struct的名字和表名必须对应，当一个struct的数据要放入两个不同名结构一样的表中，就非常麻烦。
而这里会提供灵活的数据结构，用户修改后可实现自己的目的。
不做成傻瓜化的操作，暴露给最外层sql语句，也会助于用户更容易调试，优化。不会让用户完全不知道内部生成的sql是什么样子，是否用到索引的困惑
*/

/*
创建，一般放在init里初始化，全项目统一使用一个
var conf = "gechengzhen:123456@tcp(172.16.1.61:3306)/userdata?charset=utf8"
p := mysqlx.DbInit( conf,3 )
fmt.Println( p.Query("select * from ytk_car_test limit 1 ", nil ))
fmt.Println( p.Exec("insert into ytk_car_test set license = '赣B'", nil))


如果是多sql事务
tx := db.Begin()
tx.Query(xxx)
tx.Commit()


query的结果是QueryRes ，本质是一个map，可以批量修改，然后QueryRes 可以转化成Struct业务模型,可单个也可以批量。


*/

//代表事务
type DbTx struct {
	realtx *sql.Tx
}

//提交事务
func (tx *DbTx) Commit() error {
	return tx.realtx.Commit()
}

//事务回滚
func (tx *DbTx) Rollback() error {
	return tx.realtx.Rollback()
}

//代表数据库操作，自带池子
type DbPool struct {
	realPool *sql.DB
}

func DbInit(conStr string, maxOpenConns int) (*DbPool, error) {
	if maxOpenConns <= 0 || maxOpenConns >= 1000 {
		maxOpenConns = 20
	}
	pool, err := sql.Open("mysql", conStr)
	if err == nil {
		pool.SetMaxOpenConns(maxOpenConns)
		pool.SetMaxIdleConns(maxOpenConns / 5)
		pool.SetConnMaxLifetime(time.Second * 1200)
		p := &DbPool{}
		p.realPool = pool
		return p, nil
	}
	return nil, err
}

//获取*sql.DB
func (p *DbPool) DB() *sql.DB {
	return p.realPool
}

func (p *DbPool) Begin() (*DbTx, error) {
	realtx, err := p.realPool.Begin()
	if err != nil {
		return nil, err
	} else {
		t := &DbTx{}
		t.realtx = realtx
		return t, nil
	}
}
