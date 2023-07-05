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
	api_shorten_bulk "github.com/jei-el/vuo.be-backend/src/api/shorten-bulk"
	config "github.com/jei-el/vuo.be-backend/src/config"
)

func Serve() {
	r := chi.NewRouter()

	api_shorten_bulk.NewShortenBulkModule().Init(r)

	port := os.Getenv(config.SERVER_PORT)
	if port == "" {
		log.Fatalf("Server port not found")
	}

	log.Printf("Start server at port %s with env '%s'.", port, os.Getenv(config.CURRENT_ENV))

	port = ":" + port
	http.ListenAndServe(port, r)
}
