package dashboard

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var (
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
	prefix := getHttpPrefix()
	routes := []route{
		newPage(prefix+"/", indexTmpls, getIndexData),
	}

	router := mux.NewRouter().StrictSlash(true)
	for _, r := range routes {
		log.Printf("Registering route for %q on %q\n", r.Method(), r.Pattern())
		router.
			Methods(r.Method()).
			Path(r.Pattern()).
			HandlerFunc(r.HandlerFunc())
	}
	return router
}

func getHttpPrefix() string {
	return os.Getenv("DASHBOARD_HTTP_PREFIX")
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
			log.Printf("error getting template data: %v\n", err)
			serveISE(w)
			return
		}
		err = p.tmpl.ExecuteTemplate(w, baseTemplate, data)
		if err != nil {
			log.Printf("error rendering template: %v\n", err)
			serveISE(w)
			return
		}
	}

	log.Printf("Auth is disabled is set, not checking credentials\n")
	return fn
}
