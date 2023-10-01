package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/julienschmidt/httprouter"
	"snippetbox.anukuljoshi/internals/models"
)

// struct to hold form data and field errors
type snippetCreateForm struct {
	Title string
	Content string
	Expires int
	FieldErrors map[string]string
}

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
	data := app.newTemplateData(r)
	// initialize snippetCreateForm struct to pass to template
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.tmpl.html", data)
}

// handler for creating a snippet
func (app *application) createSnippetPost(w http.ResponseWriter, r *http.Request) {
	// call r.ParseForm() to parse and add POST request body data in r.PostForm map
	err := r.ParseForm()
	if err!=nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// use Get method on r.PostForm to get POST data
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	// r.PostForm.Get() return string
	// convert string to int
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err!=nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// initialize instance of snippetCreateForm to hold form data and empty fields errors
	form := snippetCreateForm{
		Title: title,
		Content: content,
		Expires: expires,
		FieldErrors: make(map[string]string),
	}
	// validations check for title
	// 1. title is not empty
	// 2. title is less than 100 characters
	if strings.TrimSpace(form.Title)=="" {
		form.FieldErrors["title"] = "This field cannot be blank"
	}else if utf8.RuneCountInString(form.Title)>100 {
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}
	// validations check for content
	// 1. content is not empty
	if strings.TrimSpace(form.Content)=="" {
		form.FieldErrors["content"] = "This field cannot be blank"
	}
	// validation checks for expires
	// expires should be either 1, 7 or 365
	if form.Expires!=1 && form.Expires!=7 && form.Expires!=365 {
		form.FieldErrors["form.Expires"] = "This field must be equal to 1, 7 or 365"
	}
	// return bad request if form.FieldErrors are present
	if len(form.FieldErrors)>0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusBadRequest, "create.tmpl.html", data)
		return
	}
	// call insert for snippet model with data
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err!=nil {
		app.serverError(w, err)
		return
	}
	// redirect to snippet view for the created snippet id
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
