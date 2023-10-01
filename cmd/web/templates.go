package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"snippetbox.anukuljoshi/internals/models"
)

type templateData struct {
	CurrentYear int
	Snippet *models.Snippet
	Snippets []*models.Snippet
	Form any
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

// convert time.Time to human readable format
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 03:04PM")
}

// initialize template.FuncMap object and store in global variable
// lookup table for template function and our created functions
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// initialize map to act as cache
	cache := map[string]*template.Template{}
	// get a slice of all path which match the pattern
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err!=nil {
		return nil, err
	}
	for _, page := range pages {
		// extract filename from path
		name := filepath.Base(page)
		// register template.FuncMap before calling ParseFiles() method
		// create empty template set and register function with Funcs method
		// parse base template into a template set
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
		if err!=nil {
			return nil, err
		}
		// call ParseGlob() on base template set to add all partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err!=nil {
			return nil, err
		}
		// call ParseFiles() to add page to the template set
		ts, err = ts.ParseFiles(page)
		if err!=nil {
			return nil, err
		}
		// add template set to cache with filename as key
		cache[name] = ts
	}
	return cache, nil
}
