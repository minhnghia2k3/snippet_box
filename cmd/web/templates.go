package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/minhnghia2k3/snippet_box/internal/models"
)

// Holding structure for any dynamic data
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	Form        any
	Flash       string
}

// Function humanDate() which is returns a nicely formatted string
// representation of time.Time object
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// Initialize a template.FuncMap object and store it in a global variable.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Get the name of the file
		name := filepath.Base(page)

		// The template.FuncMap must be registered with the template set
		// before call ParseFiles().
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() *on this template set* to add any partials.
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		// Call ParseFile() *on this template set* to add the page template.
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map.
		// cache[home] = *Template
		cache[name] = ts
	}
	return cache, nil
}
