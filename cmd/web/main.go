package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/minhnghia2k3/snippet_box/internal/models"
)

type config struct {
	addr      string
	staticDir string
	dsn       string
}

var (
	Red   = "\033[31m"
	Blue  = "\033[34m"
	White = "\033[97m"
	Gray  = "\033[37m"
)

// struct application will inject to the handlers.
type application struct {
	debug          bool
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       models.SnippetModelInterface
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

/*
==========================================================================================
							MAIN APPLICATION
==========================================================================================
*/

/*
This function will:

- Parsing command-line flag values from user stdin.

- Defining log level.

- Opening database connection.

- Initializing template cache.

- Initializing application struct for inject to another handlers.

- Creating a server instance.
*/
func main() {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	debug := flag.Bool("debug", false, "Application debug mode")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static address")
	flag.StringVar(&cfg.dsn, "dsn", "web:secret@tcp(localhost:3306)/snippetbox?parseTime=true", "MySQL data source name")
	// Must call before use the addr variable
	flag.Parse()

	infoLog := log.New(os.Stdout, Blue+"[INFO]\t"+Gray, log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, Red+"[ERROR]\t"+Gray, log.Lshortfile|log.Ldate|log.Ltime)

	db, err := openDB(cfg.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	// Initialize a new instance of our application struct
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}
	formDecoder := form.NewDecoder()
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	// Serve request on HTTPS only
	sessionManager.Cookie.Secure = true

	app := &application{
		debug:          *debug,
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:      cfg.addr,
		Handler:   app.routes(),
		ErrorLog:  errorLog,
		TLSConfig: tlsConfig,

		// Add Idle, Read and Write timeouts to the server.
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Server is listening on port %s", cfg.addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

// The openDB() function wraps sql.Open and returns a sql.DB connection pool.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
