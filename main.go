package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// handler for catch all
func home(w http.ResponseWriter,  r *http.Request){
	// return NotFound if url does not match "/" exactly
	if r.URL.Path!="/"{
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello World from snippetbox"))
}

// handler for viewing a snippet
func viewSnippet(w http.ResponseWriter,  r *http.Request){
	// get id sent in query params
	id, err := strconv.Atoi(r.URL.Query().Get("id"))	
	if err != nil || id < 1{
		http.NotFound(w, r)
		return
	}
	// use FprintF to write formatted string to response writer
	fmt.Fprintf(w, "display a specific snippet with ID %d", id)
}

// handler for viewing a snippet
func createSnippet(w http.ResponseWriter,  r *http.Request){
	// check if  request method is POST
	if r.Method != "POST"{
		w.Header().Set("Allow", http.MethodPost)
		// use http.Error to send 405 status code and "Method Not Allowed" message
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return	
	}
	w.Write([]byte("create a new snippet"))
}

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