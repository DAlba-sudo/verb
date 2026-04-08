package htmx

import (
	"html/template"
	"net/http"

	"strings"
)

const (
	hxWrapper = `
	{{- define "attrs" -}}
		{{- if .Htmx.Class }} class="{{ .Htmx.Class }}" {{ end -}}
		style="padding: 0; margin: 0;"
		{{- if and .Htmx.HxAjax (ne .Htmx.HxAjax.Method "none") }} hx-{{ .Htmx.HxAjax.Method }}="{{ .Htmx.HxAjax.URL }}"{{ end -}}
		{{- if .Htmx.HxTrigger }} hx-trigger="{{ .Htmx.HxTrigger }}"{{ end -}}
		{{- if .Htmx.HxTarget }} hx-target="{{ .Htmx.HxTarget }}"{{ end -}}
		{{- if .Htmx.HxSwap }} hx-swap="{{ .Htmx.HxSwap }}"{{ end -}}
		{{- if .Htmx.HxInclude }} hx-include="{{ .Htmx.HxInclude }}"{{ end -}}
	{{- end -}}

	{{- if eq .Htmx.HxContainerTag "span" -}}
		<span {{ template "attrs" . }}>{{ template "content" . }}</span>
	{{- else if eq .Htmx.HxContainerTag "li" -}}
		<li {{ template "attrs" . }}>{{ template "content" . }}</li>
	{{- else if eq .Htmx.HxContainerTag "tr" -}}
		<tr {{ template "attrs" . }}>{{ template "content" . }}</tr>
	{{- else if eq .Htmx.HxContainerTag "td" -}}
		<td {{ template "attrs" . }}>{{ template "content" . }}</td>
	{{- else -}}
		<div {{ template "attrs" . }}>{{ template "content" . }}</div>
	{{- end -}}
	`

	MethodGET    = "get"
	MethodPOST   = "post"
	MethodPUT    = "put"
	MethodPATCH  = "patch"
	MethodDELETE = "delete"
	MethodNone   = "none"
)

type ajax struct {
	Method string
	URL    string
}

type Htmx struct {
	HxContainerTag string
	HxAjax         ajax
	HxTrigger      string
	HxTarget       string
	HxSwap         string
	HxInclude      string
	Class          string
}

func (hx Htmx) Data(w http.ResponseWriter, r *http.Request) (any, error) {
	return hx, nil
}

func (hx Htmx) Name() string {
	return "Htmx"
}

func (h *Htmx) GET(url string) *Htmx {
	h.HxAjax = ajax{Method: MethodGET, URL: url}
	return h
}

func (h *Htmx) POST(url string) *Htmx {
	h.HxAjax = ajax{Method: MethodPOST, URL: url}
	return h
}

func (h *Htmx) PUT(url string) *Htmx {
	h.HxAjax = ajax{Method: MethodPUT, URL: url}
	return h
}

func (h *Htmx) PATCH(url string) *Htmx {
	h.HxAjax = ajax{Method: MethodPATCH, URL: url}
	return h
}

func (h *Htmx) DELETE(url string) *Htmx {
	h.HxAjax = ajax{Method: MethodDELETE, URL: url}
	return h
}

func (h *Htmx) Trigger(trigger ...string) *Htmx {
	h.HxTrigger = strings.Join(trigger, ",")
	return h
}

func (h *Htmx) Target(target string) *Htmx {
	h.HxTarget = target
	return h
}

func (h *Htmx) Tag(tag string) *Htmx {
	h.HxContainerTag = tag
	return h
}

func (h *Htmx) Swap(swap string) *Htmx {
	h.HxSwap = swap
	return h
}

func (h *Htmx) Include(include string) *Htmx {
	h.HxInclude = include
	return h
}

func (h *Htmx) Classes(classes ...string) *Htmx {
	h.Class = strings.Join(classes, " ")
	return h
}

func (h *Htmx) Build(content string, funcs map[string]any) *template.Template {
	// We define the base "htmx" template using the wrapper
	// and then define the "content" template with the user's input.
	t, err := template.New("htmx").Parse(hxWrapper)
	if err != nil {
		panic(err)
	}

	if funcs != nil {
		_, err = t.New("content").Funcs(funcs).Parse(content)
	} else {
		_, err = t.New("content").Parse(content)
	}
	if err != nil {
		panic(err)
	}

	return t
}

func Div() *Htmx {
	return &Htmx{
		HxAjax:         ajax{Method: MethodNone, URL: ""},
		HxContainerTag: "div",
	}
}

func Create(tag string) *Htmx {
	return &Htmx{
		HxAjax:         ajax{Method: MethodNone, URL: ""},
		HxContainerTag: tag,
	}
}
