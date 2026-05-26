package verb

import (
	"errors"
	"html/template"
	"net/http"
	"os"
	"strings"
)

var (
	ErrRouteNotFound    = errors.New("route not found")
	ErrTemplateNotFound = errors.New("template not found")
)

type Server struct {
	Map   map[string](map[string]*Route)
	Funcs template.FuncMap

	Start ServerStartOptions
}

func CreateServer() *Server {
	s := &Server{
		Map:   make(map[string](map[string]*Route)),
		Funcs: make(template.FuncMap),
	}

	return s
}

type ServerStartOptions struct {
	Reload bool
}

func (s *Server) Register(method string, url string, templatePath string) (*Route, error) {
	route := &Route{}

	// read content from file
	data, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, err
	}

	// build the template
	tmpl, err := template.New("content").Funcs(s.Funcs).Parse(string(data))
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

func (s *Server) defaultHandler(route *Route, w http.ResponseWriter, r *http.Request) error {
	if route.Template == nil {
		return ErrTemplateNotFound
	}

	model := map[string]any{}
	for _, bridge := range route.Bridges {
		data, key, err := bridge(w, r)
		if err != nil {
			continue
		}

		model[key] = data
	}

	return route.Template.ExecuteTemplate(w, "content", model)
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

	return http.ListenAndServe(address, mux)
}
