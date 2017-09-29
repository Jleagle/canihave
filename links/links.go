package links

import (
	"net/http"
	"net/url"
)

func GetHeaderLinks(r *http.Request) (links map[string]string) {

	links = map[string]string{
		"sortDate":  makeSearchLink(r, "sort", "date"),
		"sortPop":   makeSearchLink(r, "sort", "pop"),
		"sortRank":  makeSearchLink(r, "sort", "rank"),
		"sortPrice": makeSearchLink(r, "sort", "price"),
	}

	return links
}

func makeSearchLink(r *http.Request, override ...string) (link string) {

	query := r.URL.Query()
	params := url.Values{}

	sort := query.Get("sort")
	if sort != "" {
		params.Add("sort", sort)
	}

	category := query.Get("category")
	if category != "" {
		params.Add("category", category)
	}

	search := query.Get("search")
	if search != "" {
		params.Add("search", search)
	}

	page := query.Get("page")
	if page != "" {
		params.Add("page", page)
	}

	if len(override) == 1 {
		params.Del(override[0])
	} else if len(override) == 2 {
		params.Set(override[0], override[1])
	}

	return "/?" + params.Encode()
}
