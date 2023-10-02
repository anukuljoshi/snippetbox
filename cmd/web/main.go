package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"snippetbox.anukuljoshi/internals/models"
)

// Define an application struct to hold the application-wide dependencies
type application struct {
	debug bool
	errorLog *log.Logger
	infoLog *log.Logger
	snippets *models.SnippetModel
	users *models.UserModel
	templateCache map[string]*template.Template
	formDecoder *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	// define a new command line flag "addr" to specify to host address
	addr := flag.String("addr", ":4000", "HTTP network address")
	debug := flag.Bool("debug", false, "Enable debug mode")
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

	// initialize a template cache
	templateCache, err := newTemplateCache()
	if err!=nil {
		errorLog.Fatal(err)
	}

	// initialize form decoder instance
	formDecoder := form.NewDecoder()

	// initialize new SessionManager
	var sessionManager = scs.New()
	// use mysql as db
	sessionManager.Store = mysqlstore.New(db)
	// set Lifetime to 12 hours, session automatically expires after 12 hours
	sessionManager.Lifetime = 12 * time.Hour
	// set Secure attribute to send Cookie only in https connection
	sessionManager.Cookie.Secure = true

	// initialize an application struct with dependencies
	app := &application{
		debug: *debug,
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &models.SnippetModel{DB: db},
		users: &models.UserModel{DB: db},
		templateCache: templateCache,
		formDecoder: formDecoder,
		sessionManager: sessionManager,
	}

	// initialize a tls.Config struct to hold non-default TLS settings
	var tlsConfig = &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}
	// initialize a new http.Server struct
	server := &http.Server{
		Addr: *addr,
		Handler: app.routes(),
		ErrorLog: errorLog,
		TLSConfig: tlsConfig,
		// add Idle, Read and Write timeouts
		IdleTimeout: time.Minute,
		ReadTimeout: 5*time.Second,
		WriteTimeout: 10*time.Second,
	}
	infoLog.Printf("Starting server on %s", *addr)
	err = server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
