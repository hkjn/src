package dashboard

import (
	"html/template"
	"net/http"

	"hkjn.me/src/googleauth"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

var (
	authDisabled = false
	baseTmpls    = []string{
		"tmpl/base.tmpl",
		"tmpl/scripts.tmpl",
		"tmpl/style.tmpl",
	}
	indexTmpls = append(
		baseTmpls,
		"tmpl/index.tmpl",
		"tmpl/links.tmpl",
		"tmpl/prober.tmpl",
	)
	baseTemplate = "base"
)

// newRouter returns a new router for the endpoints of the dashboard.
//
// newRouter panics if the config wasn't loaded.
func newRouter(debug bool) *mux.Router {
	routes := []route{
		newPage("/", indexTmpls, getIndexData),
		simpleRoute{"/connect", "GET", googleauth.ConnectHandler},
	}
	if debug {
		authDisabled = true // TODO(hkjn): Avoid global variable.
	}

	router := mux.NewRouter().StrictSlash(true)
	for _, r := range routes {
		glog.V(1).Infof("Registering route for %q on %q\n", r.Method(), r.Pattern())
		router.
			Methods(r.Method()).
			Path(r.Pattern()).
			HandlerFunc(r.HandlerFunc())
	}
	return router
}

// getTemplate returns the template loaded from the paths.
//
// getTemplate parses the .tmpl files from disk.
func getTemplate(tmpls []string) *template.Template {
	return template.Must(template.ParseFiles(tmpls...))
}

// serveISE serves an internal server error to the user.
func serveISE(w http.ResponseWriter) {
	http.Error(w, "Internal server error.", http.StatusInternalServerError)
}

// route describes how to serve HTTP on an endpoint.
type route interface {
	Method() string                // GET, POST, PUT, etc.
	Pattern() string               // URI for the route
	HandlerFunc() http.HandlerFunc // HTTP handler func
}

// simpleRoute implements the route interface for endpoints.
type simpleRoute struct {
	pattern, method string
	handlerFunc     http.HandlerFunc
}

func (r simpleRoute) Method() string { return r.method }

func (r simpleRoute) Pattern() string { return r.pattern }

func (r simpleRoute) HandlerFunc() http.HandlerFunc { return r.handlerFunc }

// getDataFn is a function to get template data.
type getDataFn func(http.ResponseWriter, *http.Request) (interface{}, error)

// page implements the route interface for endpoints that render HTML.
type page struct {
	pattern         string
	tmpl            *template.Template // backing template
	getTemplateData getDataFn
}

// newPage returns a new page.
func newPage(pattern string, tmpls []string, getData getDataFn) *page {
	return &page{
		pattern,
		getTemplate(tmpls),
		getData,
	}
}

func (p page) Method() string { return "GET" }

func (p page) Pattern() string { return p.pattern }

// HandlerFunc returns the http handler func, which renders the
// template with the data.
func (p page) HandlerFunc() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		data, err := p.getTemplateData(w, r)
		if err != nil {
			glog.Errorf("error getting template data: %v\n", err)
			serveISE(w)
			return
		}
		err = p.tmpl.ExecuteTemplate(w, baseTemplate, data)
		if err != nil {
			glog.Errorf("error rendering template: %v\n", err)
			serveISE(w)
			return
		}
	}

	if authDisabled {
		glog.V(1).Infof("Auth is disabled is set, not checking credentials\n")
	} else {
		fn = googleauth.RequireLogin(fn)
	}
	return fn
}
