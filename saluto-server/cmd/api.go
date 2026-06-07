package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Fylo4/saluto/tutorial"
	"github.com/jackc/pgx/v5"

	"github.com/caarlos0/env/v11"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

type EnvConfig struct {
	User    string
	Pass    string
	Host    string
	Port    string
	Name    string
	SSLMode string
}

func (app *application) mount() (http.Handler, error) {
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

	ctx := context.Background()

	// Load environment variables
	envName := os.Getenv("APP_ENV")
	if envName == "" {
		envName = "local"
	}

	err := godotenv.Load(".env." + envName)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		return nil, err
	}

	envVars := EnvConfig{
		User:    os.Getenv("DB_USER"),
		Pass:    os.Getenv("DB_PASS"),
		Host:    os.Getenv("DB_HOST"),
		Port:    os.Getenv("DB_PORT"),
		Name:    os.Getenv("DB_NAME"),
		SSLMode: os.Getenv("DB_SSLMODE"),
	}
	if err := env.Parse(&envVars); err != nil {
		log.Fatalf("Error reading the environment variables: %v", err)
		return nil, err
	}

	var tlsConfig *tls.Config

	if envName == "local" {
		// Load AWS RDS CA certificate
		rootCertPool := x509.NewCertPool()
		pem, err := ioutil.ReadFile("../aws-certificate-bundle.pem")
		if err != nil {
			log.Fatal("failed to read CA file:", err)
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			log.Fatal("failed to append CA cert")
		}

		// Build TLS config
		tlsConfig = &tls.Config{
			RootCAs:    rootCertPool,
			ServerName: envVars.Host,
		}
	}

	// Build pgx config
	connstr := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		envVars.User, envVars.Pass, envVars.Host, envVars.Port, envVars.Name, envVars.SSLMode,
	)
	log.Println(connstr)
	config, err := pgx.ParseConfig(connstr)
	if err != nil {
		panic(err)
	}
	if tlsConfig != nil {
		config.TLSConfig = tlsConfig
	}

	// Connect
	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		log.Fatal("connection failed:", err)
	}

	log.Println("Connected")
	defer conn.Close(ctx)
	sql := tutorial.New(conn)

	r.Get("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{\"message\": \"Hello world\"}"))
	})
	r.Get("/api/writers", func(w http.ResponseWriter, r *http.Request) {
		_, err := sql.ListAuthors(ctx)
		if err != nil {
			log.Fatal(err)
		}
		w.Write([]byte("{\"message\": \"TODO Implement me\"}"))
	})

	return r, nil
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
