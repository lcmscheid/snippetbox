package main

import (
	"database/sql"
	"flag"
	"github.com/lcmscheid/snippetbox/pkg/models/mysql"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Addr      string
	StaticDir string
	dsn       string
}

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	config        *Config
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	cfg := new(Config)
	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&cfg.dsn, "dsn", "web:manager@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	app := &application{
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		config:   cfg,
	}

	db, err := openDb(cfg.dsn)
	if err != nil {
		app.errorLog.Fatal(err)
	}
	defer db.Close()
	app.snippets = &mysql.SnippetModel{DB: db}

	templateCache, err := newTemplateCache("./ui/html")
	if err != nil {
		app.errorLog.Fatal(err)
	}
	app.templateCache = templateCache

	srv := &http.Server{
		Addr:     app.config.Addr,
		ErrorLog: app.errorLog,
		Handler:  app.routes(),
	}

	app.infoLog.Printf("Starting server on %s", cfg.Addr)
	err = srv.ListenAndServe()
	app.errorLog.Fatal(err)
}

func openDb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
