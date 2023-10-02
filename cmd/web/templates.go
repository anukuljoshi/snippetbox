package main

import (
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"time"

	"github.com/justinas/nosurf"
	"snippetbox.anukuljoshi/internals/models"
	"snippetbox.anukuljoshi/ui"
)

type templateData struct {
	CurrentYear int
	Snippet *models.Snippet
	Snippets []*models.Snippet
	User *models.User
	Form any
	Flash any
	IsAuthenticated bool
	CSRFToken string
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
		// add flash message if it exists
		Flash: app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken: nosurf.Token(r),
	}
}

// convert time.Time to human readable format
func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
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
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl.html")
	if err!=nil {
		return nil, err
	}
	for _, page := range pages {
		// extract filename from path
		name := filepath.Base(page)
		// create a slice containing filepath for templates to parse
		patterns := []string{
			"html/base.tmpl.html",
			"html/partials/*.tmpl.html",
			page,
		}
		if err!=nil {
			return nil, err
		}
		// call ParseFiles() to add page to the template set
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err!=nil {
			return nil, err
		}
		// add template set to cache with filename as key
		cache[name] = ts
	}
	return cache, nil
}
