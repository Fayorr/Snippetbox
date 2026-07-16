package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql" // New import
	"github.com/joho/godotenv"
	"snippetbox.fayokunmiosho.com/internal/models"
)

type application struct {
	logger         *slog.Logger
	snippets models.SnippetModelInterface 
	users models.UserModelInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}


func openDB(dsn string) (*sql.DB, error) {

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}
	err = db.Ping()

	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	_ = godotenv.Load()

	dbUser := getEnv("DB_USER", "web")
	dbPass := getEnv("DB_PASSWORD", "pass")
	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "snippetbox")
	
	addr := flag.String("addr", getEnv("ADDR", ":4000"), "HTTP network address")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		// AddSource: true,
	}))

	db, err := openDB(dsn)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
		users:          &models.UserModel{DB: db},
	}
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.CurveP256, tls.X25519},
		MinVersion:       tls.VersionTLS13,
	}
	srv := &http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("Starting Server at PORT", slog.String("addr", *addr))

	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")

	app.logger.Error(err.Error())
	os.Exit(1)

}
