package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Jleagle/canihave/models"
	"github.com/Jleagle/canihave/store"
	"github.com/Masterminds/squirrel"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {

	search := r.Form.Get("search")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))

	vars := searchVars{}
	vars.Items = handleQuery(page, search)
	vars.Search = search
	vars.Search64 = base64.StdEncoding.EncodeToString([]byte(search))
	vars.Javascript = []string{"/assets/search.js", "//platform.twitter.com/widgets.js"}

	returnTemplate(w, "search", vars)
}

func ajaxHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))

	vars := searchVars{}
	vars.Items = handleQuery(page, query.Get("search"))

	returnTemplate(w, "search_ajax", vars)
}

func handleQuery(page int, search string) []models.Item {

	if page < 1 {
		page = 1
	}

	// Make SQL
	conn := store.GetMysqlConnection()
	query := squirrel.Select("*").From("items")
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	fmt.Printf("%v", (page-1)*12)
	query = query.OrderBy("dateCreated DESC").Limit(12).Offset(uint64((page - 1) * 12))

	sql, args, error := query.ToSql()
	if error != nil {
		fmt.Println(error)
	}
	fmt.Printf("%v", sql)

	// Run SQL
	rows, error := conn.Query(sql, args...)
	if error != nil {
		fmt.Println(error)
	}
	defer rows.Close()

	// Convert to types
	results := []models.Item{}
	item := models.Item{}
	for rows.Next() {
		rows.Scan(&item.ID, &item.DateCreated, &item.DateUpdated, &item.Name, &item.Desc, &item.Source)
		results = append(results, item)
	}

	return results
}

func s2uint64(i string) uint64 {

	ix, err := strconv.ParseInt(i, 10, 64)
	if err != nil {
		panic(err)
	}
	return uint64(ix)
}

type searchVars struct {
	Items      []models.Item
	Page       string
	Search     string
	Search64   string
	Javascript []string
}
