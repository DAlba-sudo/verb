package verb

import (
	"errors"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrRouteNotFound    = errors.New("route not found")
	ErrTemplateNotFound = errors.New("template not found")
)

type Server struct {
	Map     map[string](map[string]*Route)
	Funcs   template.FuncMap
	Options ServerOptions
}

type ServerOptions struct {
	StaticFilesDir  string
	ReloadTemplates bool
}

func CreateServer(static string) *Server {
	s := &Server{
		Map:   make(map[string](map[string]*Route)),
		Funcs: make(template.FuncMap),
		Options: ServerOptions{
			StaticFilesDir: static,
		},
	}

	return s
}

func (s *Server) Func(key string, fn any) {
	s.Funcs[key] = fn
}

func filesToTemplates(files []string, funcs template.FuncMap) (*template.Template, error) {
	var tmpl *template.Template

	// read content from file
	for _, templatePath := range files {
		data, err := os.ReadFile(templatePath)
		if err != nil {
			return nil, err
		}

		name := filepath.Base(templatePath)

		// build the template
		if tmpl == nil {
			tmpl, err = template.New(name).Funcs(funcs).Parse(string(data))
		} else {
			_, err = tmpl.New(name).Funcs(funcs).Parse(string(data))
		}
		if err != nil {
			return nil, err
		}
	}

	return tmpl, nil
}

func (s *Server) Register(method string, url string, templatePaths ...string) (*Route, error) {
	route := &Route{
		Metadata: routeMetadata{
			OriginalTemplatePaths: templatePaths,
		},
	}

	tmpl, err := filesToTemplates(templatePaths, s.Funcs)
	if err != nil {
		return nil, err
	}

	route.Template = tmpl

	// register route
	endpoints, ok := s.Map[url]
	if !ok {
		s.Map[url] = make(map[string]*Route)
		endpoints = s.Map[url]
	}
	endpoints[method] = route

	return route, nil
}

func (s *Server) API(method string, url string, handler func(w http.ResponseWriter, r *http.Request) error) (*Route, error) {
	route := &Route{
		Handler: handler,
	}

	endpoints, ok := s.Map[url]
	if !ok {
		s.Map[url] = make(map[string]*Route)
		endpoints = s.Map[url]
	}
	endpoints[method] = route

	return route, nil
}

func (s *Server) defaultHandler(route *Route, w http.ResponseWriter, r *http.Request) error {
	if route.Template == nil {
		return ErrTemplateNotFound
	}

	if s.Options.ReloadTemplates {
		tmpl, err := filesToTemplates(route.Metadata.OriginalTemplatePaths, s.Funcs)
		if err == nil {
			route.Template = tmpl
		}
	}

	model := map[string]any{}
	for _, bridge := range route.Bridges {
		data, key, err := bridge(w, r)
		if err != nil {
			continue
		}

		model[key] = data
	}

	return route.Template.Execute(w, model)
}

func (s Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	path := normalizeURL(r.URL.Path)

	routes, ok := s.Map[path]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	route, ok := routes[r.Method]
	if !ok {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if route.Handler == nil && route.Template != nil {
		s.defaultHandler(route, w, r)
		return
	}

	route.Handler(w, r)
}

func normalizeURL(url string) string {
	return strings.ToLower(strings.TrimSpace(strings.TrimSuffix(url, "/")))
}

func (s *Server) Serve(address string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleRequest)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(s.Options.StaticFilesDir))))

	return http.ListenAndServe(address, mux)
}
