package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"github.com/katarzynakawala/diffs-app/pkg/models"
	"github.com/katarzynakawala/diffs-app/pkg/models/mysql"
)

type contextKey string

const contextKeyIsAuthenticated = contextKey("isAuthenticated")

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	snippets      interface {
		Insert(string, string, string) (int, error)
		Get(int) (*models.Snippet, error)
		Latest() ([]*models.Snippet, error)
	}
	templateCache map[string]*template.Template
	users         interface {
		Insert(string, string, string)  error
		Authenticate(string, string) (int, error)
		Get(int) (*models.User, error)
	}
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dns", "web:pass@/diffsapp?parseTime=true", "MySQL data source name")
	secret := flag.String("secret", "f7Mfh+gFbnvHgS*+5Ph8qFWhTzbpd@rt", "Secret key")
	flag.Parse()

	//logger for writting information messages
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	//logger for writting error messages
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	//new session
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	defer db.Close()

	//new template cache
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	//new instance of application with dependencies
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
		users: &mysql.UserModel{DB: db},
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	//creation of struct to use custom logger with every error
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

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
