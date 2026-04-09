package verb

import (
	"net/http"
)

// A bridge is used to populate templates. It requires two functions, one
// for specifying the data to be passed to the template, and another for specifying
// the name.
type Bridge interface {
	Data(http.ResponseWriter, *http.Request, map[string]any) (any, error)
	Name() string
}

type DataBridge struct {
	Key     string
	Provide func(r *http.Request, model map[string]any) (any, error)
}

func (d DataBridge) Name() string {
	return d.Key
}

func (d DataBridge) Data(w http.ResponseWriter, r *http.Request, model map[string]any) (any, error) {
	return d.Provide(r, model)
}

func Map(name string, provider func(*http.Request, map[string]any) (any, error)) DataBridge {

	return DataBridge{
		Key:     name,
		Provide: provider,
	}
}
