package render

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
)

type HTMLRenderer struct {
	templates map[string]*template.Template
}

func NewHTMLRenderer(templateFS fs.FS, templatesFiles ...string) (*HTMLRenderer, error) {
	cache := make(map[string]*template.Template)

	pages, err := fs.Glob(templateFS, "*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		if name == "base.html" {
			continue
		}

		// Each page template is compiled with base.html so define "main" is isolated
		tmpl, err := template.ParseFS(templateFS, "base.html", name)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", name, err)
		}
		cache[name] = tmpl
	}

	return &HTMLRenderer{
		templates: cache,
	}, nil
}

func (hr *HTMLRenderer) Render(w http.ResponseWriter, status int, name string, data any) {
	var body bytes.Buffer
	if err := hr.Execute(&body, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	_, _ = body.WriteTo(w)
}

func (hr *HTMLRenderer) Execute(w io.Writer, name string, data any) error {
	tmpl, ok := hr.templates[name]
	if !ok {
		return fmt.Errorf("template %q not found", name)
	}
	return tmpl.ExecuteTemplate(w, "base", data)
}
