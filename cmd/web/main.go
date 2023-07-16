package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// Define an application struct to hold the application-wide dependencies
type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
}

func main() {
	// define a new command line flag "addr" to specify to host address
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// create a new logger for info messages
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// create a new logger for info messages
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// initialize an application struct with dependencies
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
	}
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

	// initialize a new http.Server struct
	server := &http.Server{
		Addr: *addr,
		Handler: mux,
		ErrorLog: errorLog,
	}
	// use http.ListenAndServe to create a new web server
	// pass in two parameters
	// 1. TCP network address
	// 2. ServeMux created earlier
	infoLog.Printf("Starting server on %s", *addr)
	err := server.ListenAndServe()
	errorLog.Fatal(err)
}