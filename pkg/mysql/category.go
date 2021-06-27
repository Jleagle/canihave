package mysql

import (
	logger2 "github.com/Jleagle/canihave/pkg/logger"
	"github.com/Masterminds/squirrel"
)

type Category struct {
	ID       string
	Category string
	Count    string
}

func (c Category) GetLink() (link string) {
	return "/?category=" + c.ID
}

//todo, make a struct and add these methods as listeners
// add a get function that just fills in the restof the struct

func CategoryIDFromName(name string) (id string) {

	// Make the query
	query := squirrel.Select("id").From("categories").Where(squirrel.Eq{"name": name}).Limit(1)
	err := QueryRow(query).Scan(&id)
	if err != nil {
		logger2.Err("Can't scan category row", err)
	}

	return id
}

func CategoryNameFromID(id string) (name string) {

	// Make the query
	builder := squirrel.Select("amazon").From("categories").Where(squirrel.Eq{"id": id}).Limit(1)
	err := QueryRow(builder).Scan(&name)
	if err != nil {
		logger2.Err("Can't scan category row", err)
	}

	return name
}

func SaveCategory(name string) (success bool, err error) {

	builder := squirrel.Insert("categories").Columns("amazon").Values(name)
	err = Insert(builder)
	if err == nil {
		return true, err
	}

	logger2.Err("Can't insert category", err)
	return false, err
}
