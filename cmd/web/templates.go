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
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
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
		// parse base template into a template set
		ts, err := template.ParseFiles("./ui/html/base.tmpl.html")
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
