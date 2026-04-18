# (Verb) Interfaces

The following document outlines the interfaces available to interact with the 
verb framework.

## Data Structures

The following data structures are used in the verb framework. 

#### HTMX 

This HTMX data structure is used to return a dynamic HTMX HTML container. You can interact with it using the following 
methods:

```go
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

func (h *Htmx) Vals(vals any) *Htmx {
	data, err := json.Marshal(vals)
	if err != nil {
		return h
	}

	h.HxVals = string(data)
	return h
}

func (h *Htmx) SelfEncodeRequest() *Htmx {
	h.HxRedoEncode = true
	return h
}
```

- This structure abides by the Bridge interface described below, so it can be used as a bridge to expose data to the template when rendering a page. This allows for dynamic HTMX containers to be rendered on the page with data from the server.

#### Route 

A route contains endpoint specific information:

- The type of route (Page, Component, Action), this is used in the HTTP handler to multiplex the request to the correct handler.
- URL, this is used to match the incoming request to the correct route.
- Bridges, this is a slice of "Bridge" interfaces that are registered to the route and return an object and error combo that is 
used passed to the template when rendering the page.
- Error, this is a slice of "Bridge" interfaces that are run when a bridge returns an error. This is used to handle errors gracefully and return a custom error page instead of a 500 error.
- Miscellaneous objects like the `html/template` object, and an `htmx` object are also in this data structure.

#### Bridge (Interface)

A bridge exposes data to the template when rendering a page. It must conform to the following interface:

```go
// A bridge is used to populate templates. It requires two functions, one
// for specifying the data to be passed to the template, and another for specifying
// the name.
type Bridge interface {
	Data(http.ResponseWriter, *http.Request, map[string]any) (any, error)
	Name() string
}
```

- The `Name()` function's return value is used as the key with which the template can access the data returned by the `Data(...)` function.
- The map[string]any parameter in the `Data(...)` function is used to pass data from one bridge to another. This allows for bridges to build on top of each other and create a data pipeline that can be used to populate the template with complex data structures, or to not repeat work.

## Methods 

### Page

```go
func (v *Verb) Page(url string, file string) *Route
```

A Page exposes an HTML file's contents on a given URL. It registers the URL so it can 
only be reached from a GET request to that route.

- The HTML contents are embedded into a base html template file.
- A user can register a "Bridge" with the "Route" to expose data to the template via go's html/template package.

### Component

```go
func (v *Verb) Component(file string, hx *htmx.Htmx) *Route 
```

A "Component" takes an HTML file whose name is used to craft the HTMX endpoint it can be retrieved 
from. The format pre-pends `/htmx/` to the file name, and removes the file extension. For example, if the file name is `user.html`, then the component can be retrieved from the URL `/htmx/user`. 

- Path components are ignored (i.e., a file at `/components/user/form.html` would be available at `/htmx/form`).

### Action

```
func (v *Verb) Action(method string, url string, handler func(http.ResponseWriter, *http.Request) error) *Route 
```

An action functions similarly to a component, except the route is not auto-registered and a custom handler can be used 
as well as the bridge structure. This is particularly useful for returning HTMX responses to POST form data.
