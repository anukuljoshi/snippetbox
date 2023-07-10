package main

import (
	"log"
	"net/http"
)

func main() {
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
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", viewSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// use http.ListenAndServe to create a new web server
	// pass in two parameters
	// 1. TCP network address
	// 2. ServeMux created earlier
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}