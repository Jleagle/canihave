package helpers

import "net/http"

func GetQueryParam(r *http.Request, key string, defaultValue string) string {
	v := r.URL.Query().Get(key)
	if v != "" {
		v = defaultValue
	}
	return v
}
