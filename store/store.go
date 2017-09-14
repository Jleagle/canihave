package store

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/Masterminds/squirrel"
	"github.com/bradfitz/gomemcache/memcache"
)

var mysqlConnection *sql.DB
var memcacheConnection *memcache.Client
var mysqlPrepareStatements map[string]*sql.Stmt

func GetMemcacheConnection() *memcache.Client {

	if memcacheConnection == nil {
		memcacheConnection = memcache.New("127.0.0.1:11211")
	}
	return memcacheConnection
}

func GetMysqlConnection() *sql.DB {

	if mysqlConnection == nil {

		password := os.Getenv("SQL_PW")
		if len(password) > 0 {
			password = ":" + password
		}

		var err error
		mysqlConnection, err = sql.Open("mysql", "root"+password+"@tcp(127.0.0.1:3306)/canihave")
		if err != nil {
			panic(err.Error())
		}
	}

	return mysqlConnection
}

func getPrepareStatement(query string) (statement *sql.Stmt) {

	if mysqlPrepareStatements == nil {
		mysqlPrepareStatements = make(map[string]*sql.Stmt)
	}

	byteArray := md5.Sum([]byte(query))
	hash := hex.EncodeToString(byteArray[:])

	if val, ok := mysqlPrepareStatements[hash]; ok {
		if ok {
			return val
		}
	}

	conn := GetMysqlConnection()

	var err error
	statement, err = conn.Prepare(query)
	if err != nil {
		panic(err.Error())
	}

	mysqlPrepareStatements[hash] = statement
	return statement
}

func Insert(builder squirrel.InsertBuilder) (rows *sql.Rows, err error) {

	rawSQL, args, err := builder.ToSql()
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(rawSQL)

	prep := getPrepareStatement(rawSQL)

	rows, err = prep.Query(args...)

	return rows, err
}

func Query(builder squirrel.SelectBuilder) (rows *sql.Rows) {

	rawSQL, args, err := builder.ToSql()
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(rawSQL)

	prep := getPrepareStatement(rawSQL)

	rows, err = prep.Query(args...)
	if err != nil {
		fmt.Println(err)
	}

	return rows
}

func QueryRow(builder squirrel.SelectBuilder) *sql.Row {

	rawSQL, args, err := builder.ToSql()
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(rawSQL)

	prep := getPrepareStatement(rawSQL)

	return prep.QueryRow(args...)
}
