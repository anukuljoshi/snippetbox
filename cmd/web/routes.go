package main

import "net/http"

// returns a servemux containing our application routes
func (app *application) routes() http.Handler {
	// create a new ServeMux
	// register home function as handler for "/" path
	mux := http.NewServeMux()

	// create a file server for serving static files
	fileServer := http.FileServer(http.Dir("./ui/static"))
	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// other application routes
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.viewSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	return secureHeaders(mux)
}
