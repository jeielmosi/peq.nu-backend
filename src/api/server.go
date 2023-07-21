package api

import (
	/*
	   "net/http"

	   "github.com/go-chi/chi/v5/middleware"
	*/
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	api_shorten_bulk "github.com/jeielmosi/peq.nu-backend/src/api/shorten-bulk"
	config "github.com/jeielmosi/peq.nu-backend/src/config"
)

func Serve() {
	r := chi.NewRouter()
	r.Use(middleware.DefaultLogger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	api_shorten_bulk.NewShortenBulkModule().Init(r)

	port := os.Getenv(config.PORT)
	if port == "" {
		log.Fatalf("Server port not found\n")
	}

	log.Printf("Start peq.nu serverless at port %s\n", port)

	port = ":" + port
	http.ListenAndServe(port, r)
}
