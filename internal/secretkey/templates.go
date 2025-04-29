package secretkey

import (
	"fmt"
	"bytes"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

type templateData struct {
	Error     string
	Success   string
}

func humaData(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humaData,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {

	cache := make(map[string]*template.Template)

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))

	if err != nil {
		return nil, err
	}
	for _, page := range pages {

		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}
		cache[name] = ts

	}

	return cache, nil
}

func render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := templateCache[name]
	if !ok {
		serverError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, nil)
	if err != nil {
		serverError(w, err)
		return
	}
	buf.WriteTo(w)
}

