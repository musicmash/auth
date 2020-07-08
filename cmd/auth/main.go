package main

import (
	"flag"
	"log"
	"time"

	"github.com/go-chi/chi"
	"github.com/musicmash/auth/internal/api/handlers/auth"
	"github.com/musicmash/auth/internal/api/router"
	"github.com/musicmash/auth/internal/api/server"
	"github.com/musicmash/auth/internal/backends/firebase"
)

func main() {
	// define cli args
	serviceAccountFilePath := flag.String(
		"path-to-service-account",
		"/etc/auth/serviceAccountKey.json",
		"abs path to serviceAccountKey.json",
	)
	flag.Parse()

	// make backend
	backend, err := firebase.New(*serviceAccountFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// make router
	authHandler := auth.NewHandler(backend)
	r := chi.NewRouter()
	r.Post("/auth", authHandler.DoAuth)
	router := router.New(r)

	// make http server
	server := server.New(router, &server.Options{
		IP:           "0.0.0.0",
		Port:         1200,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	})

	// and finally listen
	log.Fatal(server.ListenAndServe())
}
