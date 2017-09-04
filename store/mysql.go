package store

import (
	"database/sql"
	"os"
)

var mysql *sql.DB
var mysqlInsertItem *sql.Stmt

func GetMysqlConnection() *sql.DB {

	if mysql == nil {

		password := os.Getenv("SQL_PW")
		if len(password) > 0 {
			password = ":" + password
		}

		var err error
		mysql, err = sql.Open("mysql", "root"+password+"@tcp(127.0.0.1:3306)/canihave")
		if err != nil {
			panic(err.Error())
		}
	}

	return mysql
}

func GetInsertPrep() *sql.Stmt {

	if mysqlInsertItem == nil {

		conn := GetMysqlConnection()

		var err error
		mysqlInsertItem, err = conn.Prepare("INSERT INTO items (id, dateCreated, dateUpdated, name, link, source, salesRank, photo, productGroup, price, region) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			panic(err.Error())
		}
	}

	return mysqlInsertItem
}
