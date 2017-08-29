package store

import (
	"database/sql"
	"os"
)

var mysql *sql.DB

func GetMysqlConnection() *sql.DB {

	if mysql == nil {

		password := os.Getenv("SQL_PW")
		if len(password) > 0 {
			password = ":" + password
		}

		var error error
		mysql, error = sql.Open("mysql", "root"+password+"@tcp(127.0.0.1:3306)/canihave")
		if error != nil {
			panic(error.Error())
		}
	}

	return mysql
}
