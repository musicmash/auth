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
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
}

type SpotifyApplication struct {
	ID     string `yaml:"id"`
	Secret string `yaml:"secret"`
}
