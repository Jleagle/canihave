package models

import (
	"github.com/Jleagle/canihave/logger"
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
	err := store.QueryRow(query).Scan(&id)
	if err != nil {
		logger.Err("Can't scan category row", err)
	}

	return id
}

func CategoryNameFromID(id string) (name string) {

	// Make the query
	builder := squirrel.Select("amazon").From("categories").Where(squirrel.Eq{"id": id}).Limit(1)
	err := store.QueryRow(builder).Scan(&name)
	if err != nil {
		logger.Err("Can't scan category row", err)
	}

	return name
}

func SaveCategory(name string) (success bool, err error) {

	builder := squirrel.Insert("categories").Columns("amazon").Values(name)
	err = store.Insert(builder)
	if err == nil {
		return true, err
	}

	logger.Err("Can't insert category", err)
	return false, err
}
