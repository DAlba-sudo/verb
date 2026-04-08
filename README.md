# Verb Framework README

Verb is a Go web framework built specifically for developers who want to leverage HTMX to create dynamic, server-driven user interfaces with minimal friction. It streamlines the process of mapping HTML templates to routes while providing a robust bridge system for data injection.

### Defining Pages

Pages in Verb are intended to be the primary entry points for your application. When you define a page using the `Page` method, the framework reads your HTML file and embeds it into a base layout known as `base.html`. This allows you to maintain a consistent structure, such as headers and footers, across your entire site. The content of your page file is treated as a template definition named "content" which is then executed within the context of the base template.

For example, your `base.html` might look like this:
```html
<html>
  <body>
    <main>{{ template "content" . }}</main>
  </body>
</html>
```
When you register a route with `v.Page("/", "home.html")`, the contents of `home.html` will automatically fill that "content" slot.

### Creating HTMX Components

Components offer a more granular way to build interfaces and are ideal for the partial updates that HTMX excels at. Unlike pages, components do not pull from the `base.html` template. Instead, they are wrapped in a specific HTML container tag—such as a `div`, `span`, `tr`, or `li`—and enriched with HTMX attributes directly from your Go code.

When you define a component like `v.Component("user-card.html", htmx.Div().GET("/api/user"))`, the framework generates a specific URL for it, typically prefixed with `/htmx/`. The resulting HTML delivered to the browser will be wrapped in the specified tag with the corresponding HTMX attributes:

```html
<div hx-get="/api/user" style="padding: 0; margin: 0;">
  </div>
```

You can chain various methods in Go to configure your component's behavior, such as `Target`, `Swap`, `Trigger`, and `Include`, which map directly to `hx-target`, `hx-swap`, `hx-trigger`, and `hx-include`.

### Data Management with Bridges

Bridges are the mechanism for passing data from your Go logic into your HTML templates. A Bridge associates a specific key with a provider function that fetches data during a request. You can register global bridges in the server settings to act as middleware for every route, or attach them to specific routes for localized data.

Inside your HTML, the data returned by the bridge is accessible through the key you provided. For instance, if you have a bridge with the key "Profile", you can access its fields in your template using dot notation, such as `{{ .Profile.Username }}`.

### Custom Template Functions

To enhance your templates further, you can register custom Go functions that will be available globally across all pages and components. By using the `Func` method, you can add helpers for formatting dates, manipulating strings, or any other logic you prefer to handle at the template level. These functions are injected into the template's function map before parsing, ensuring they are ready for use in your HTML.
