package mysql

import (
	"github.com/Masterminds/squirrel"
)

type Category struct {
	ID       string
	Category string
	Count    string
}

func GetCategories(offset uint64) (categories []Category, err error) {

	builder := squirrel.Select("*")
	builder = builder.From("categories")
	builder = builder.OrderBy("id ASC")
	builder = builder.Limit(offset)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	c, err := getClient()
	if err != nil {
		return nil, err
	}

	err = c.Select(&categories, sql, args...)
	return categories, err
}

type CategoryCount struct {
	Node  string
	Count int
}

func GetCategoryCounts(typex string) (counts []CategoryCount, err error) {

	builder := squirrel.Select("nodeName AS category, count(nodeName) AS count")
	builder = builder.From("items")
	builder = builder.GroupBy("nodeName")
	builder = builder.OrderBy("count DESC")
	builder = builder.Where("type = ?", TypeScraper)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	c, err := getClient()
	if err != nil {
		return nil, err
	}

	err = c.Select(&counts, sql, args...)
	return counts, err
}
