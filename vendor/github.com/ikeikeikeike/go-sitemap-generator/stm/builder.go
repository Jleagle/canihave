package stm

import (
	"fmt"
)

var poolBuffer = NewBufferPool()

// BuilderError provides interface for it can confirm the error in some difference.
type BuilderError interface {
	error
	FullError() bool
}

// Builder provides interface for adds some kind of url sitemap.
type Builder interface {
	XMLContent() []byte
	Content() []byte
	Add(interface{}) BuilderError
	Write()
}

// SitemapURL provides generated xml interface.
type SitemapURL interface {
	XML() []byte
}

// Attrs defines for xml attribute.
type Attrs []interface{}

// Attr defines for xml attribute.
type Attr map[string]string

// URL User should use this typedef in main func.
type URL map[string]interface{}

// URLJoinBy that's convenient.
func (u URL) URLJoinBy(key string, joins ...string) URL {
	var values []string
	for _, k := range joins {
		values = append(values, fmt.Sprint(u[k]))
	}

	u[key] = URLJoin("", values...)
	return u
}

// BungURLJoinBy that's convenient. Though, this is Bung method.
func (u *URL) BungURLJoinBy(key string, joins ...string) {
	orig := *u

	var values []string
	for _, k := range joins {
		values = append(values, fmt.Sprint(orig[k]))
	}

	orig[key] = URLJoin("", values...)
	*u = orig
}

// type News map[string]interface{}
