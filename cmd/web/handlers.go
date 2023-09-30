package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"snippetbox.anukuljoshi/internals/models"
)

// handler for catch all
func (app *application) home(w http.ResponseWriter,  r *http.Request){
	snippets, err := app.snippets.Latest()
	if err!=nil {
		app.serverError(w, err)
		return
	}
	// call newTemplateData to create templateData with CurrentYear
	data := app.newTemplateData(r)
	data.Snippets = snippets
	// use the render helper method
	app.render(w, http.StatusOK, "home.tmpl.html", data)
}

// handler for viewing a snippet
func (app *application) viewSnippet(w http.ResponseWriter,  r *http.Request){
	params := httprouter.ParamsFromContext(r.Context())
	// get id sent in query params
	id, err := strconv.Atoi(params.ByName("id"))	
	if err != nil || id < 1{
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err!=nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		}
		app.serverError(w, err)
	}
	// call newTemplateData to create templateData with CurrentYear
	data := app.newTemplateData(r)
	data.Snippet = snippet
	// use render helper method
	app.render(w, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display snippet form for creating new snippet"))
}

// handler for creating a snippet
func (app *application) createSnippetPost(w http.ResponseWriter, r *http.Request) {
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n-Kobayashi Issa"
	expires := 7
	// call insert for snippet model with data
	id, err := app.snippets.Insert(title, content, expires)
	if err!=nil {
		app.serverError(w, err)
		return
	}
	// redirect to snippet view for the created snippet id
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
