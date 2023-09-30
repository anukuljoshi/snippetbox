package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"snippetbox.anukuljoshi/internals/models"
)

// Define an application struct to hold the application-wide dependencies
type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
	snippets *models.SnippetModel
}

func main() {
	// define a new command line flag "addr" to specify to host address
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// create a new logger for info messages
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// create a new logger for info messages
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	err := godotenv.Load(".env")
	if err!=nil {
		errorLog.Fatal(err)
	}
	dsn := os.Getenv("MYSQL_DSN")

	db, err := openDB(dsn)
	if err!=nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// initialize an application struct with dependencies
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &models.SnippetModel{DB: db},
	}
	// initialize a new http.Server struct
	server := &http.Server{
		Addr: *addr,
		Handler: app.routes(),
		ErrorLog: errorLog,
	}
	infoLog.Printf("Starting server on %s", *addr)
	err = server.ListenAndServe()
	errorLog.Fatal(err)
}

// The openDB() function wraps sql.Open() and 
// returns a sql.DB connection pool for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err!=nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
