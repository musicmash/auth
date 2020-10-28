package config

import "time"

type AppConfig struct {
	HTTP               HTTPConfig         `yaml:"http"`
	DB                 DBConfig           `yaml:"db"`
	Log                LogConfig          `yaml:"log"`
	SpotifyApplication SpotifyApplication `yaml:"spotify"`
}

type HTTPConfig struct {
	IP           string        `yaml:"ip"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
	DomainName   string        `yaml:"domain_name"`
}

type LogConfig struct {
	Level string `yaml:"level"`
}

type DBConfig struct {
	Host                  string        `yaml:"host"`
	Port                  int           `yaml:"port"`
	Name                  string        `yaml:"name"`
	Login                 string        `yaml:"login"`
	Password              string        `yaml:"password"`
	MaxOpenConnections    int           `yaml:"max_open_connections"`
	MaxIdleConnections    int           `yaml:"max_idle_connections"`
	MaxConnectionLifeTime time.Duration `yaml:"max_connection_life_time"`
	MaxConnectionIdleTime time.Duration `yaml:"max_connection_idle_time"`
}

type SpotifyApplication struct {
	ID     string `yaml:"id"`
	Secret string `yaml:"secret"`
}
