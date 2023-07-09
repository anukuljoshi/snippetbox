package main

import (
	"log"
	"net/http"
)

func main() {
	// create a new ServeMux
	// register home function as handler for "/" path
	mux := http.NewServeMux()
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