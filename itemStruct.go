package main

import (
	"fmt"
	"time"

	cache "github.com/patrickmn/go-cache"
)

// item is the database row
type item struct {
	ID          string
	DateUpdated string
	DateCreated string
	TimesAdded  string
	Name        string
	Desc        string
	Source      string
}

func (i item) GetUKLink() string {
	return "https://www.amazon.co.uk/dp/" + i.ID + "?tag=canihaveone00-21"
}

func (i item) GetUKPixel() string {
	return "//ir-uk.amazon-adsystem.com/e/ir?t=canihaveone00-21&l=am2&o=2&a=" + i.ID
}

func (i item) inMemcache() bool {
	c := cache.New(5*time.Minute, 10*time.Minute)
	_, found := c.Get("foo")
	return found
}

func (i item) inMysql() bool {

	db := connectToSQL()

	item := item{}
	error := db.QueryRow("SELECT id FROM items WHERE id = ? LIMIT 1", i.ID).Scan(&item.ID)
	if error != nil {
		fmt.Println(error)
	}

	return true
}

func (i item) saveToMemcache() {
	c := cache.New(5*time.Minute, 10*time.Minute)
	c.Set("foo", "bar", cache.DefaultExpiration)
}

func (i item) getFromMemcache() item {
	c := cache.New(5*time.Minute, 10*time.Minute)
	foo, found := c.Get("foo")
	if found {
		return foo.(item)
	}
	return i.getFromMysql()
}

func (i item) getFromMysql() item {

	return i
}

func (i item) saveToMysql() {

	i.saveToMemcache()
}

func (i item) refresh() item {

	//get from amazon save into mysql and memcache

	return i
}
