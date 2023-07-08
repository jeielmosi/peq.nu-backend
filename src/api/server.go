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
	"github.com/go-chi/cors"
	api_shorten_bulk "github.com/jei-el/vuo.be-backend/src/api/shorten-bulk"
	config "github.com/jei-el/vuo.be-backend/src/config"
)

func Serve() {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	api_shorten_bulk.NewShortenBulkModule().Init(r)

	port := os.Getenv(config.PORT)
	if port == "" {
		log.Fatalf("Server port not found")
	}

	log.Printf("Start server at port %s!!!", port)

	port = ":" + port
	http.ListenAndServe(port, r)
}
