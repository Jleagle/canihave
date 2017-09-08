package store

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Masterminds/squirrel"
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
		mysqlInsertItem, err = conn.Prepare("INSERT INTO items (id, dateCreated, dateUpdated, name, link, source, salesRank, photo, productGroup, price, region, hits, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			panic(err.Error())
		}
	}

	return mysqlInsertItem
}

func QueryRows(queryBuilder squirrel.SelectBuilder) *sql.Rows {

	rawSQL, args, err := queryBuilder.ToSql()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(rawSQL)

	// Run SQL
	rows, err := GetMysqlConnection().Query(rawSQL, args...)
	if err != nil {
		fmt.Println(err)
	}

	return rows
}

func QueryRow(queryBuilder squirrel.SelectBuilder) *sql.Row {

	rawSQL, args, err := queryBuilder.ToSql()
	if err != nil {
		fmt.Println(err)
	}

	return GetMysqlConnection().QueryRow(rawSQL, args...)
}
