package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	handler "github.com/musicmash/auth/internal/api/handlers/spotify"
	"github.com/musicmash/auth/internal/api/server"
	"github.com/musicmash/auth/internal/backend"
	"github.com/musicmash/auth/internal/log"
)

func main() {
	const callbackPath = "/v1/spotify/auth-callback"

	// parse cli args
	domainName := flag.String("nginx-domain-name", os.Getenv("NGINX_DOMAIN_NAME"), "domain name for building redirect url")
	spotifyAppID := flag.String("spotify-app-id", os.Getenv("SPOTIFY_ID"), "spotify application client id")
	spotifyAppSecret := flag.String("spotify-app-secret", os.Getenv("SPOTIFY_SECRET"), "spotify application secret key")
	flag.Parse()

	// validate cli args
	if len(*domainName) == 0 {
		exitIfError(errors.New("nginx domain name is empty, so we can't build redirect url"))
	}

	if len(*spotifyAppID) == 0 || len(*spotifyAppSecret) == 0 {
		exitIfError(errors.New("spotify application credentials are empty"))
	}

	redirectURL := fmt.Sprintf("https://%s%s", *domainName, callbackPath)
	b := backend.New(redirectURL, *spotifyAppID, *spotifyAppSecret)
	h := handler.NewHandler(b)

	// setup logger
	log.SetWriters(log.GetConsoleWriter())

	// make router
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Get(callbackPath, h.DoAuth)

	// make http server
	server := server.New(r, &server.Options{
		IP:           "0.0.0.0",
		Port:         1200,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	})

	log.Infof("Please log in to Spotify by visiting the following page in your browser: %s", b.GetAuthURL("auth"))

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
