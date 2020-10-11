package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/musicmash/auth/internal/log"
	yaml "gopkg.in/yaml.v2"
)

func New() *AppConfig {
	return &AppConfig{
		HTTP: HTTPConfig{
			IP:           "0.0.0.0",
			Port:         1200,
			DomainName:   "musicmash.me",
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  10 * time.Second,
		},
		DB: DBConfig{
			Host:     "musicmash.db",
			Port:     5432,
			Name:     "auth",
			Login:    "auth",
			Password: "auth_pass",
		},
		Log: LogConfig{
			Level: "INFO",
		},
		SpotifyApplication: SpotifyApplication{
			ID:     "",
			Secret: "",
		},
	}
}

func (c *AppConfig) LoadFromFile(configPath string) error {
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	return c.LoadFromBytes(b)
}

func (c *AppConfig) LoadFromBytes(val []byte) error {
	if err := yaml.Unmarshal(val, c); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (c *AppConfig) FlagSet() {
	flag.StringVar(&c.HTTP.IP, "http-ip", c.HTTP.IP, "API ip address")
	flag.IntVar(&c.HTTP.Port, "http-port", c.HTTP.Port, "API port")
	flag.StringVar(&c.HTTP.DomainName, "http-domain-name", c.HTTP.DomainName, "server domain name for spotify callback")

	flag.StringVar(&c.DB.Host, "db-host", c.DB.Host, "Database host")
	flag.IntVar(&c.DB.Port, "db-port", c.DB.Port, "Database port")
	flag.StringVar(&c.DB.Name, "db-name", c.DB.Name, "Database name")
	flag.StringVar(&c.DB.Login, "db-login", c.DB.Login, "Database user login")
	flag.StringVar(&c.DB.Password, "db-pass", c.DB.Password, "Database user password")

	flag.StringVar(&c.Log.Level, "log-level", c.Log.Level, "Log level")

	flag.StringVar(&c.SpotifyApplication.ID, "spotify-app-id", c.SpotifyApplication.ID, "spotify application client id")
	flag.StringVar(&c.SpotifyApplication.Secret, "spotify-app-secret", c.SpotifyApplication.Secret, "spotify application secret key")
}

func (c *AppConfig) FlagReload() {
	_ = flag.Set("http-ip", c.HTTP.IP)
	_ = flag.Set("http-port", fmt.Sprintf("%d", c.HTTP.Port))
	_ = flag.Set("http-domain-name", c.HTTP.DomainName)

	_ = flag.Set("db-host", c.DB.Host)
	_ = flag.Set("db-port", fmt.Sprint(c.DB.Port))
	_ = flag.Set("db-name", c.DB.Name)
	_ = flag.Set("db-login", c.DB.Login)
	_ = flag.Set("db-pass", c.DB.Password)

	_ = flag.Set("log-level", c.Log.Level)

	_ = flag.Set("spotify-app-app-id", c.SpotifyApplication.ID)
	_ = flag.Set("spotify-app-secret", c.SpotifyApplication.Secret)
}

func (c *AppConfig) Dump() string {
	b, _ := yaml.Marshal(c)
	return string(b)
}

func (db *DBConfig) GetConnString() string {
	return fmt.Sprintf(
		"host=%v port=%v user=%v dbname=%v sslmode=disable password=%v",
		db.Host, db.Port, db.Login, db.Name, db.Password)
}
