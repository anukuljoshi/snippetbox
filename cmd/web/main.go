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
	// initialize a new http.Server struct
	server := &http.Server{
		Addr: *addr,
		Handler: app.routes(),
		ErrorLog: errorLog,
	}
	infoLog.Printf("Starting server on %s", *addr)
	err := server.ListenAndServe()
	errorLog.Fatal(err)
}