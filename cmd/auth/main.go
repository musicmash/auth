package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/musicmash/auth/internal/api/handlers/auth"
	"github.com/musicmash/auth/internal/api/handlers/spotify"
	"github.com/musicmash/auth/internal/api/server"
	"github.com/musicmash/auth/internal/backend"
	"github.com/musicmash/auth/internal/config"
	"github.com/musicmash/auth/internal/db"
	"github.com/musicmash/auth/internal/log"
	"github.com/musicmash/auth/internal/version"
)

func main() {
	_ = flag.Bool("version", false, "Show build info and exit")
	if versionRequired() {
		_, _ = fmt.Fprintln(os.Stdout, version.FullInfo)
		os.Exit(0)
	}

	// parse conf
	conf := config.New()
	conf.FlagSet()
	configPath := flag.String("config", "", "Path to auth.yml config")
	_ = flag.Bool("help", false, "Show this message and exit")
	if helpRequired() {
		flag.PrintDefaults()
		os.Exit(0)
	}

	flag.Parse()
	if *configPath != "" {
		if err := conf.LoadFromFile(*configPath); err != nil {
			exitIfError(err)
		}

		// set not provided flags as config values
		conf.FlagReload()

		// override config values with provided flags
		flag.Parse()
	}
	exitIfError(validateConfig(conf))

	mgr, err := db.New(conf.DB.GetConnString())
	exitIfError(err)

	const callbackPath = "/v1/spotify/auth-callback"
	redirectURL := fmt.Sprintf("https://%s%s", conf.HTTP.DomainName, callbackPath)
	b := backend.New(mgr, redirectURL, conf.SpotifyApplication.ID, conf.SpotifyApplication.Secret)
	spotifyCallbackHandler := spotify.NewHandler(b)
	authHandler := auth.NewHandler(mgr)

	// setup logger
	log.SetWriters(log.GetConsoleWriter())

	// make router
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Get(callbackPath, spotifyCallbackHandler.DoAuth)
	r.Post("/auth", authHandler.DoAuth)

	// make http server
	server := server.New(r, &server.Options{
		IP:           conf.HTTP.IP,
		Port:         conf.HTTP.Port,
		ReadTimeout:  conf.HTTP.ReadTimeout,
		WriteTimeout: conf.HTTP.WriteTimeout,
		IdleTimeout:  conf.HTTP.IdleTimeout,
	})

	log.Infof("Please log in to Spotify by visiting the following page in your browser: %s", b.GetAuthURL("auth"))

	// and finally listen
	exitIfError(server.ListenAndServe())
}

func validateConfig(conf *config.AppConfig) error {
	if conf.Log.Level == "" {
		conf.Log.Level = "info"
	}

	if conf.HTTP.Port < 0 || conf.HTTP.Port > 65535 {
		return errors.New("invalid port: value should be in range: 0 < value < 65535")
	}

	if len(conf.HTTP.DomainName) == 0 {
		return errors.New("nginx domain name is empty, so we can't build redirect url")
	}

	if len(conf.SpotifyApplication.ID) == 0 || len(conf.SpotifyApplication.Secret) == 0 {
		return errors.New("spotify application credentials are empty")
	}

	return nil
}

func isArgProvided(argName string) bool {
	for _, arg := range os.Args {
		if strings.Contains(arg, argName) {
			return true
		}
	}
	return false
}

func helpRequired() bool {
	return isArgProvided("-help")
}

func versionRequired() bool {
	return isArgProvided("-version")
}

func exitIfError(err error) {
	if err == nil {
		return
	}

	log.Error(err.Error())
	os.Exit(2)
}
