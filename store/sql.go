package store

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/Jleagle/canihave/logger"
	"github.com/Masterminds/squirrel"
)

var mysqlConnection *sql.DB

var mysqlPrepareStatements map[string]*sql.Stmt

func Query(builder squirrel.SelectBuilder) (rows *sql.Rows) {

	rawSQL, args, err := builder.ToSql()
	if err != nil {
		logger.Err("Can't make query SQL", err)
	}

	//logger.Info("SQL: " + rawSQL)

	prep := getPrepareStatement(rawSQL)

	rows, err = prep.Query(args...)
	if err != nil {
		logger.Err("Can't query prepped statement", err)
	}

	return rows
}

func QueryRow(builder squirrel.SelectBuilder) *sql.Row {

	rawSQL, args, err := builder.ToSql()
	if err != nil {
		logger.Err("Can't make query SQL", err)
	}

	//logger.Info("SQL: " + rawSQL)

	prep := getPrepareStatement(rawSQL)

	return prep.QueryRow(args...)
}

func Insert(builder squirrel.InsertBuilder) (err error) {

	rawSQL, args, err := builder.ToSql()
	if err != nil {
		logger.Err("Can't make insert SQL", err)
	}

	//logger.Info("SQL: " + rawSQL)

	prep := getPrepareStatement(rawSQL)

	_, err = prep.Exec(args...)

	return err
}

func Update(builder squirrel.UpdateBuilder) (err error) {

	rawSQL, args, err := builder.ToSql()
	if err != nil {
		logger.Err("Can't make update SQL", err)
	}

	//logger.Info("SQL: " + rawSQL)

	prep := getPrepareStatement(rawSQL)

	_, err = prep.Exec(args...)

	return err
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

	conn := getMysqlConnection()

	var err error
	statement, err = conn.Prepare(query)
	if err != nil {
		logger.Err("Can't run prepared statement", err)
	}

	mysqlPrepareStatements[hash] = statement
	return statement
}

func getMysqlConnection() *sql.DB {

	if mysqlConnection == nil {

		database := os.Getenv("CANIHAVE_MYSQL_DATABASE")
		username := os.Getenv("CANIHAVE_MYSQL_USERNAME")
		password := os.Getenv("CANIHAVE_MYSQL_PASSWORD")
		if len(password) > 0 {
			password = ":" + password
		}

		dsn := fmt.Sprintf("%s%s@tcp(%s:%s)/%s",
			username, password, "127.0.0.1", "3306", database)

		var err error
		mysqlConnection, err = sql.Open("mysql", dsn)
		if err != nil {
			logger.Err("Can not connect to MySQL", err)
		}
	}

	return mysqlConnection
}
