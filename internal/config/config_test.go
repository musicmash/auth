package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig_LoadFromFile(t *testing.T) {
	// arrange
	assert.NoError(t, os.Setenv("DB_HOST", "auth.db"))
	assert.NoError(t, os.Setenv("DB_PORT", "5432"))
	assert.NoError(t, os.Setenv("DB_NAME", "auth"))
	assert.NoError(t, os.Setenv("DB_USER", "auth"))
	assert.NoError(t, os.Setenv("DB_PASSWORD", "auth"))
	assert.NoError(t, os.Setenv("SPOTIFY_CLIENT_ID", "2c7a0f0a-29fe-4ec4-926f-1e956297af9e"))
	assert.NoError(t, os.Setenv("SPOTIFY_CLIENT_SECRET", "75f505b3-9e40-4d55-a693-1f2388d944dd"))
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
			Host:                  "auth.db",
			Port:                  5432,
			Name:                  "auth",
			Login:                 "auth",
			Password:              "auth",
			AutoMigrate:           true,
			MigrationsDir:         "file:///etc/auth/migrations",
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
