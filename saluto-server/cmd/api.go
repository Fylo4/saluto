package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200", "https://myapp.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(c.Handler)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{\"message\": \"Hello world\"}"))
	})

	// http.ListenAndServe(":3333", r)

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 60,
	}

	log.Printf("Server has started at address %s", app.config.addr)

	return srv.ListenAndServe()
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type application struct {
	config config
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}

type testret struct {
	message string
}
