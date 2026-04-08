package verb

import (
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/DAlba-sudo/pbf"
	"github.com/DAlba-sudo/verb/htmx"
)

const (
	routeTypeComponent = "component"
	routeTypePage      = "page"
)

// This is the basic routing component that will
// be used throughout the framework. It should hold
// everything required to render a page or component.
type Route struct {
	Type    string
	URL     string
	Bridges []Bridge

	tmpl *template.Template
}

func (r *Route) Bridge(b Bridge) *Route {
	r.Bridges = append(r.Bridges, b)
	return r
}

// A page will simply be embedded into the base template.
func (v *Verb) Page(url string, file string) *Route {

	// read the files
	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	t := template.Must(v.base.Clone())
	template.Must(t.New("content").Parse(string(data)))

	r := &Route{
		Type: routeTypePage,
		URL:  url,
		tmpl: t,
	}

	v.routes[url] = r
	v.router.Add(pbf.RouteOptions{
		Method:   http.MethodGet,
		Endpoint: url,
		Handler:  v.handle,
	})
	return r
}

// A component _can be_ an htmx item, but it doesn't have to be. The idea
// is that it won't pull from the base template and it will be exposed
// using a pre-made htmx route.
func (v *Verb) Component(file string, hx *htmx.Htmx) *Route {
	// read the file contents
	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	// we are now going to clean up the string
	file_path_parts := strings.Split(file, string(os.PathSeparator))
	last_component := strings.Trim(strings.Split(file_path_parts[len(file_path_parts)-1], ".")[0], " /\n\r\t")
	url := "/" + strings.Join([]string{"htmx", last_component}, "/")

	// this is the route object that will be used as a blueprint  to
	// perform the actual routing.
	r := &Route{
		Type: routeTypeComponent,
		URL:  url,
		tmpl: hx.Build(string(data)),
	}
	r.Bridge(hx)

	// this functionally registers the template with the routing
	// runtime.
	v.routes[url] = r
	v.router.Add(pbf.RouteOptions{
		Method:   http.MethodGet,
		Endpoint: url,
		Handler:  v.handle,
	})
	return r
}
