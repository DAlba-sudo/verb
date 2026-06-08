package verb

import (
	"html/template"
	"net/http"
)

type routeMetadata struct {
	OriginalTemplatePaths []string
}

type Route struct {
	Template *template.Template
	Handler  func(w http.ResponseWriter, r *http.Request) error
	Bridges  []func(w http.ResponseWriter, r *http.Request) (any, string, error)
	Metadata routeMetadata
}
