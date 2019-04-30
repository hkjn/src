// Support for blogs in autosite
//
// Example usage:
//   myblog := NewBlog(
//     "Some title",       // for HTML <head>
//     "blog/*/*/*.tmpl",  // pattern for posts on disk
//     "blogdomain.com",   // live domain
//     []string{           // shared templates
//       "base.tmpl",
//       "blog.tmpl",
//     },
//     []string{           // templates for special handlers
//       "base.tmpl",
//       "listing.tmpl",
//     },
//   )
//   myblog.Register()
//
// This will host pages like blogdomain.com/2014/02/Foo and
// 2016/01/Bar if there's files blog/2014/02/Foo.tmpl and
// blog/2016/01/Bar.tmpl relative to the calling package, also
// using "base.tmpl" and "blog.tmpl" for all blog posts.
//
// Special handlers for listing posts will be registered, using
// "base.tmpl" and "listing.tmpl":
//   * the index (/) will show all posts from all years.
//   * for each year YYYY that has at least one blog post, /YYYY will
//     be registered, showing all posts from that year.
//   * for each year + month YYYY/MM that have at least one blog post,
//     /YYYY/MM will be registered, showing all posts from that month
//     in that year.
//
// Within the templates for these "blog listing" handlers, the following
// data is available in addition to the usual:
//   * {{.Data.Posts}}: a slice of the posts for that unit of time (all time, one
//     year, or one month, respectively)
//   * {{.Data.TimeUnit}}: the unit of time itself ("April, 2014",
//     "2014" or "all time", respectively)
package autosite

import (
	"fmt"
	"html/template"
	"log"
	"sort"
	"strings"
)

// NewBlog creates a new blog.
//
// NewBlog panics on errors reading templates.
// TODO(hkjn): Convert to ...Options() funcs.
func NewBlog(title, glob, live string, articleTmpls []string, listingTmpls []string, logger LoggerFunc, isLive bool, tmplFuncs template.FuncMap) Site {
	b := blog{internalNew(title, glob, live, articleTmpls, logger, isLive, tmplFuncs)}
	b.addHandlers(listingTmpls)
	return &b
}

// blog is a type of autosite.
type blog struct {
	site // backing site
}

// pageData is the data needed to serve a listing page.
type pageData struct {
	TimeUnit string
	Posts    posts
}
type posts []page

// listingData is the backing data used for a listing page.
type listingData struct {
	byDate map[date]posts
	byYear map[year]posts
}

// Len gives the length of the posts.
func (p posts) Len() int {
	return len(p)
}

// Less is true when post i was published before post j.
func (p posts) Less(i, j int) bool {
	return p[i].Date.before(p[j].Date)
}

// Swap changes places of element i and j.
func (p posts) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// addListing adds a listing of blog posts.
func (b *blog) addListing(uri string, tfmt string, p posts, listingTmpls []string) {
	data := pageData{
		TimeUnit: tfmt,
		Posts:    p,
	}
	b.site.addPage(uri, date{}, data, listingTmpls)
}

// byMonth accumulates a map of year -> month -> posts.
func (b blog) listingData() listingData {
	byDate := make(map[date]posts)
	byYear := make(map[year]posts)
	for uri, p := range b.site.pages {
		parts := strings.Split(uri, "/")
		d, err := getDate(parts[1], parts[2])
		if err != nil {
			log.Fatalf("failed to parse date: %v\n", err)
		}
		byDate[d] = append(byDate[d], p)
		byYear[d.Year] = append(byYear[d.Year], p)
	}
	return listingData{byDate, byYear}
}

// addHandlers registers custom handlers for the blog.
func (b *blog) addHandlers(listingTmpls []string) {
	data := b.listingData()
	for d, inDate := range data.byDate {
		sort.Sort(sort.Reverse(inDate))
		b.addListing(fmt.Sprintf("/%d/%.2d/", d.Year, d.Month), fmt.Sprintf("%s, %d", d.Month, d.Year), inDate, listingTmpls)
	}
	for y, inYear := range data.byYear {
		sort.Sort(sort.Reverse(inYear))
		b.addListing(fmt.Sprintf("/%d/", y), fmt.Sprintf("%d", y), inYear, listingTmpls)
	}
	var allPosts posts
	for _, p := range data.byYear {
		allPosts = append(allPosts, p...)
	}
	sort.Sort(sort.Reverse(allPosts))
	b.addListing("/", "all years", allPosts, listingTmpls)
}
