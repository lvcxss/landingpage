package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"landingpage/internal/render"
	"landingpage/internal/repository"
	"landingpage/internal/repository/postgres"
	"landingpage/internal/ui"

	"github.com/joho/godotenv"
)

type config struct {
	port int
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	AdminToken string
}

type application struct {
	config   config
	logger   *log.Logger
	renderer *render.HTMLRenderer
	posts    repository.PostRepository
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")

	safePassword := url.QueryEscape(os.Getenv("POSTGRES_PASSWORD"))
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		safePassword,
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	flag.StringVar(&cfg.AdminToken, "admin-token", os.Getenv("ADMIN_TOKEN"), "Admin token for authentication")
	flag.StringVar(&cfg.db.dsn, "db-dsn", dsn, "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	renderer, err := render.NewHTMLRenderer(ui.HTMLFiles, "*.html")
	if err != nil {
		logger.Fatal("Failed to load templates: ", err)
	}

	postRepo := &postgres.PostRepo{DB: db}

	// dependency injection
	app := &application{
		config:   cfg,
		logger:   logger,
		renderer: renderer,
		posts:    postRepo,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	logger.Printf("Admin token: %s", cfg.AdminToken)
	logger.Printf("sv on port %d", cfg.port)

	err = srv.ListenAndServe()
	logger.Fatal(err)
}
