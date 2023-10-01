package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// returns a servemux containing our application routes
func (app *application) routes() http.Handler {
	// initial router
	router := httprouter.New()

	// change the default not found method for httprouter
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// create a file server for serving static files
	fileServer := http.FileServer(http.Dir("./ui/static"))
	// create route to server static files
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// create new middleware chain for dynamic routes
	var dynamic = alice.New(app.sessionManager.LoadAndSave)
	// other application routes
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.viewSnippet))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.createSnippet))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.createSnippetPost))

	// middleware chain with our standard middlewares
	// which will be used for every request
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// return standard middleware chain followed by router
	return standard.Then(router)
}
