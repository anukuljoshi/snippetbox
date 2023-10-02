package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"snippetbox.anukuljoshi/ui"
)

// returns a servemux containing our application routes
func (app *application) routes() http.Handler {
	// initial router
	router := httprouter.New()

	// change the default not found method for httprouter
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// take ui.Files embedded file system and convert it to http.FS
	// create fileserver handler with ui.Files file system
	fileServer := http.FileServer(http.FS(ui.Files))
	// create route to server static files
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	// create new middleware chain for dynamic routes
	var dynamic = alice.New(app.sessionManager.LoadAndSave, app.noSurf, app.authenticate)

	// unprotected application routes
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.viewSnippet))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignUp))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignUpPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	var protected = dynamic.Append(app.requireAuthentication)
	// protected routes
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.createSnippet))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.createSnippetPost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	// middleware chain with our standard middlewares
	// which will be used for every request
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// return standard middleware chain followed by router
	return standard.Then(router)
}
