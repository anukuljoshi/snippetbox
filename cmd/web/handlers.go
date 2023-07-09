package main

import (
	"fmt"
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
