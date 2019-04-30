// +build appengine

// ae.go holds appengine-specific bits of hkjnweb.
package hkjnweb

import (
	"net/http"

	"hkjn.me/autosite"

	"appengine"
)

func GetAeLogger(r *http.Request) autosite.Logger {
	return appengine.NewContext(r)
}

func init() {
	getLogger = GetAeLogger
	Register(!appengine.IsDevAppServer())
}
