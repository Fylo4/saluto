package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	saluto "github.com/Fylo4/saluto/sql/compiled"
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

	// Load environment variables (Local only)
	envName := os.Getenv("APP_ENV")

	if envName != "" {
		err := godotenv.Load(".env." + envName)
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
			return nil, err
		}
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
	sql := saluto.New(conn)

	r.Get("/api/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{\"message\": \"Hello world\"}"))
	})
	r.Get("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		log.Println("A")
		posts, err := sql.GetMessages(ctx)
		if err != nil {
			log.Println("An error occurred while getting messages.")
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("B")
		jsonBytes, err := json.Marshal(posts)
		if err != nil {
			http.Error(w, "Failed to encode JSON", 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonBytes)
	})
	r.Post("/api/post", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("Error reading request body"))
		}
		parsedBody := API_Post_Create_Input{}
		json.Unmarshal([]byte(body), &parsedBody)
		params := saluto.PostMessageParams{
			Displayname: parsedBody.DisplayName,
			Body:        parsedBody.Message,
		}
		err = sql.PostMessage(ctx, params)
		if err != nil {
			log.Println("DB INSERT ERROR:", err)
			http.Error(w, err.Error(), 500)
			return
		}
		w.Write([]byte("true"))
	})

	return r, nil
}

type API_Post_Create_Input struct {
	DisplayName string `json:"displayName"`
	Message     string `json:"message"`
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
