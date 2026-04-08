package verb

import (
	"net/http"
)

type DataBridge struct {
	Key     string
	Provide func(r *http.Request) (any, error)
}

func (d DataBridge) Name() string {
	return d.Key
}

func (d DataBridge) Data(w http.ResponseWriter, r *http.Request) (any, error) {
	return d.Provide(r)
}

func Map(name string, provider func(*http.Request) (any, error)) DataBridge {

	return DataBridge{
		Key:     name,
		Provide: provider,
	}
}
