package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_LoadFromFile(t *testing.T) {
	// arrange
	expected := AppConfig{
		HTTP: HTTPConfig{
			IP:         "0.0.0.0",
			Port:       1200,
			DomainName: "musicmash.me",
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
			ID:     "client_id",
			Secret: "client_secret",
		},
	}

	// action
	actual := AppConfig{}
	err := actual.LoadFromFile("../../auth.example.yml")

	// assert
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
