package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig_LoadFromFile(t *testing.T) {
	// arrange
	expected := AppConfig{
		HTTP: HTTPConfig{
			IP:           "0.0.0.0",
			Port:         1200,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  10 * time.Second,
			DomainName:   "musicmash.me",
		},
		DB: DBConfig{
			Host:                  "musicmash.db",
			Port:                  5432,
			Name:                  "auth",
			Login:                 "auth",
			Password:              "auth_pass",
			MaxOpenConnections:    10,
			MaxIdleConnections:    10,
			MaxConnectionIdleTime: 3 * time.Minute,
			MaxConnectionLifeTime: 5 * time.Minute,
		},
		Log: LogConfig{
			Level: "INFO",
		},
		SpotifyApplication: SpotifyApplication{
			ID:     "client_id",
			Secret: "client_secret",
		},
	}

	// action
	conf, err := LoadFromFile("../../auth.example.yml")

	// assert
	assert.NoError(t, err)
	assert.Equal(t, expected, *conf)
}
