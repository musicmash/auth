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
	"github.com/musicmash/auth/internal/log"
	"github.com/zmb3/spotify"
)

func main() {
	// parse cli args
	spotifyAppID := flag.String("spotify-app-id", os.Getenv("SPOTIFY_ID"), "spotify application client id")
	spotifyAppSecret := flag.String("spotify-app-secret", os.Getenv("SPOTIFY_SECRET"), "spotify application secret key")
	flag.Parse()

	if len(*spotifyAppID) == 0 || len(*spotifyAppSecret) == 0 {
		exitIfError(errors.New("spotify application credentials are empty"))
	}

	const (
		state = "auth"

		// redirectURI is the OAuth redirect URI for the application.
		// You must register an application at Spotify's developer portal
		// and enter this value.
		redirectURL = "https://dev.musicmash.me/v1/spotify/auth-callback"
	)

	auth := spotify.NewAuthenticator(redirectURL)
	auth.SetAuthInfo(*spotifyAppID, *spotifyAppSecret)
	url := auth.AuthURLWithDialog(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	h := handler.NewHandler(state, &auth)

	// setup logger
	log.SetWriters(log.GetConsoleWriter())

	// make router
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Get("/v1/spotify/auth-callback", h.DoAuth)

	// make http server
	server := server.New(r, &server.Options{
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
