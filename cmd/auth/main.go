package main

import (
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/musicmash/auth/internal/api/router"
	"github.com/musicmash/auth/internal/api/server"
	"github.com/musicmash/auth/internal/log"
)

func main() {
	// setup logger
	log.SetWriters(log.GetConsoleWriter())

	// make router
	r := chi.NewRouter()
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
	exitIfError(server.ListenAndServe())
}

func exitIfError(err error) {
	if err == nil {
		return
	}

	log.Error(err.Error())
	os.Exit(2)
}
