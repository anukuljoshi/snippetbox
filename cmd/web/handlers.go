package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// handler for catch all
func (app *application) home(w http.ResponseWriter,  r *http.Request){
	// return NotFound if url does not match "/" exactly
	if r.URL.Path!="/"{
		app.notFound(w)
		return
	}
	// read the template file
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
	}
}

// handler for viewing a snippet
func (app *application) viewSnippet(w http.ResponseWriter,  r *http.Request){
	// get id sent in query params
	id, err := strconv.Atoi(r.URL.Query().Get("id"))	
	if err != nil || id < 1{
		app.notFound(w)
		return
	}
	// use FprintF to write formatted string to response writer
	fmt.Fprintf(w, "display a specific snippet with ID %d", id)
}

// handler for viewing a snippet
func (app *application) createSnippet(w http.ResponseWriter,  r *http.Request){
	// check if  request method is POST
	if r.Method != "POST"{
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return	
	}
	w.Write([]byte("create a new snippet"))
}
