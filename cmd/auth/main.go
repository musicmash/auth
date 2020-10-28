package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/musicmash/auth/internal/api"
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

	_ = flag.Bool("help", false, "Show this message and exit")
	if helpRequired() {
		flag.PrintDefaults()
		os.Exit(0)
	}

	configPath := flag.String("config", "", "abs path to conf file")
	flag.Parse()

	if *configPath == "" {
		_, _ = fmt.Fprintln(os.Stdout, "provide abs path to config via --config argument")
		return
	}

	conf, err := config.LoadFromFile(*configPath)
	if err != nil {
		exitIfError(err)
	}
	exitIfError(validateConfig(conf))

	log.SetLevel(conf.Log.Level)
	log.SetWriters(log.GetConsoleWriter())

	log.Debug(version.FullInfo)

	log.Info("connecting to db...")
	mgr, err := db.Connect(db.Config{
		DSN:                     conf.DB.GetConnString(),
		MaxOpenConnectionsCount: conf.DB.MaxOpenConnections,
		MaxIdleConnectionsCount: conf.DB.MaxIdleConnections,
		MaxConnectionIdleTime:   conf.DB.MaxConnectionIdleTime,
		MaxConnectionLifetime:   conf.DB.MaxConnectionLifeTime,
	})
	exitIfError(err)

	log.Info("connection to the db established")

	if conf.DB.AutoMigrate {
		log.Info("applying migrations..")
		err = mgr.ApplyMigrations(conf.DB.MigrationsDir)
		if err != nil {
			exitIfError(fmt.Errorf("cant-t apply migrations: %v", err))
		}
	}

	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithTimeout(context.Background(), conf.HTTP.WriteTimeout)
	defer cancel()

	redirectURL := fmt.Sprintf("https://%s/v1/callbacks/spotify/auth\"", conf.HTTP.DomainName)
	b := backend.New(mgr, redirectURL, conf.SpotifyApplication.ID, conf.SpotifyApplication.Secret)
	router := api.GetRouter(mgr, b)
	server := api.New(router, conf.HTTP)

	go gracefulShutdown(ctx, server, quit, done)

	log.Infof("Please log in to Spotify by visiting the following page in your browser: %s", b.GetAuthURL("auth"))
	log.Infof("server is ready to handle requests at: %v", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		exitIfError(fmt.Errorf("could not listen on %v: %v", server.Addr, err))
	}

	<-done
	_ = mgr.Close()
	log.Info("auth stopped")
}

func gracefulShutdown(ctx context.Context, server *api.Server, quit <-chan os.Signal, done chan<- bool) {
	<-quit
	log.Info("server is shutting down...")

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("could not gracefully shutdown the server: %v", err)
	}
	close(done)
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
