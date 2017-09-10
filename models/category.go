package models

import (
	"fmt"
	"strconv"

	"github.com/Jleagle/canihave/store"
	"github.com/Masterminds/squirrel"
)

type Category struct {
	ID         int
	AmazonName string
}

//todo, make a struct and add these methods as listeners
// add a get function that just fills in the restof the struct

func CategoryIDFromName(name string) (id string) {

	// Make the query
	query := squirrel.Select("id").From("categories").Where(squirrel.Eq{"name": name}).Limit(1)
	sql, args, err := query.ToSql()
	if err != nil {
		fmt.Println(err)
	}

	// Run the query
	err = store.GetMysqlConnection().QueryRow(sql, args...).Scan(&id)
	if err != nil {
		fmt.Printf("%v", err.Error())
	}

	return id
}

func CategoryNameFromID(id string) (name string) {

	// Make the query
	query := squirrel.Select("amazon").From("categories").Where(squirrel.Eq{"id": id}).Limit(1)
	sql, args, err := query.ToSql()
	if err != nil {
		fmt.Println(err)
	}

	// Run the query
	err = store.GetMysqlConnection().QueryRow(sql, args...).Scan(&name)
	if err != nil {
		fmt.Printf("%v", err.Error())
	}

	return name
}

func SaveCategory(name string) (id string) {

	conn := store.GetMysqlConnection()

	res, err := conn.Exec("INSERT INTO categories (amazon) VALUES (?);", name)
	if err != nil {
		fmt.Println(err)
	}

	last, _ := res.LastInsertId()
	fmt.Println(last)
	return strconv.Itoa(int(last))
}
