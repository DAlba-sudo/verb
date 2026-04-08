package verb

import (
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/DAlba-sudo/pbf"
)

// Verb is an htmx web server framework for Go. It provides a simple
// and intuitive API for building web applications using htmx.
// Verb is designed to be easy to use and flexible,
// allowing you to build web applications quickly and efficiently.
type Verb struct {
	Settings Settings

	functions map[string]any
	routes    map[string]*Route
	router    *pbf.Router
	base      *template.Template
}

type Settings struct {
	// the path to the template directory where you are placing all your
	// templates
	Templates string

	// the path to the publicly served static files.
	Static string

	// a way to inject data into all templates, this is useful for things like user sessions, etc.
	Bridges []Bridge

	// whether the templates should be reloaded on every request,
	// useful for development.
	LiveReload bool
}

func relativeFilePath(root, path string) string {
	root = strings.TrimRight(root, "/")
	path = strings.TrimLeft(path, "/")

	return root + string(os.PathSeparator) + path
}

func New(address string, port int, s Settings) *Verb {
	data, err := os.ReadFile(relativeFilePath(s.Templates, "base.html"))
	if err != nil {
		panic(err)
	}

	r := pbf.CreateRouter()
	r.Address = address
	r.Port = port

	if s.Static != "" {
		r.Mux().HandleFunc("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(s.Static))).ServeHTTP)
	}

	return &Verb{
		Settings: s,
		router:   r,
		routes:   make(map[string]*Route),
		base:     template.Must(template.New("base").Parse(string(data))),
	}
}

// A bridge is used to populate templates. It requires two functions, one
// for specifying the data to be passed to the template, and another for specifying
// the name.
type Bridge interface {
	Data(http.ResponseWriter, *http.Request) (any, error)
	Name() string
}

func (v Verb) Serve() error {
	return v.router.Start()
}

func (v *Verb) handle(w http.ResponseWriter, r *http.Request) error {
	path := strings.TrimRight(r.URL.Path, "/")

	// the following will first check if there has already been an
	// established route.
	route, ok := v.routes[path]
	if !ok {
		http.NotFound(w, r)
		return nil
	}

	if v.Settings.LiveReload {
		data, err := os.ReadFile(relativeFilePath(v.Settings.Templates, route.originalFile))
		if err != nil {
			return err
		}

		if route.hx != nil {
			route.tmpl = route.hx.Build(string(data), v.functions)
		} else {
			t := template.Must(v.base.Clone())
			template.Must(t.New("content").Funcs(v.functions).Parse(string(data)))
			route.tmpl = t
		}
	}

	model := make(map[string]any)

	// now we are going to run the global bridges, these are
	// global to the entire app (middleware).
	for _, bridge := range v.Settings.Bridges {
		data, err := bridge.Data(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}

		if data != nil {
			model[bridge.Name()] = data
		}
	}

	// now we are going to run the route specific bridges, these
	// are local to the route.
	for _, bridge := range route.Bridges {
		data, err := bridge.Data(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}

		if data != nil {
			model[bridge.Name()] = data
		}
	}

	err := route.tmpl.Execute(w, model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}
