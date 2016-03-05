package databaseConn

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"config"
	"github.com/Centny/gwf/log"
)

var db  *sql.DB
var connConfig string = ""
var first bool = true

//set db config
func SetConnConfig(s string) {
	connConfig = s
}
//connect db
func GetConn() (*sql.DB, error) {
	if db == nil {
		var err error
		db, err = sql.Open("mysql", connConfig)
		if err != nil {
			log.E("GetConn db error--:%v", err.Error())
			return nil, err
		}
	}
	return db, nil
}
//new connect db
func GetNewConn() (*sql.DB, error) {
	if db!= nil {
		log.I("Close the old db open the new one : --:%v", connConfig)
		db.Close()
	}
	var err error
	db, err = sql.Open("mysql", connConfig)
	if err != nil {
		log.E("GetNewConn db error--:%v", err.Error())
		return nil, err
	}
	if first==true {
		first = false
		err=db.Ping()
		if nil!=err {
			log.E("connect failed: --:%v--:%v", connConfig, err);
			db=nil
			return db, err
		}
		log.D("db connect!--%v", connConfig)
	}
	return db, nil
}
//close Db
func CloseDB() {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.E("close db error--:%v", err.Error())
			return
		}
		db = nil
	}
}
//connect test db
func NewTestConn() (*sql.DB, error) {
	var conn *sql.DB
	dbConfig := config.TestDbConfig()
	SetConnConfig(dbConfig)
	log.D("your test db config is--:%v",connConfig)
	conn, _ = GetNewConn()
	return conn, nil
}
