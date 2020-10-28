package main

import "errors"

var (
	errInvalidPort      = errors.New("invalid port: value should be in range: 0 < value < 65535")
	errEmptyDomain      = errors.New("nginx domain name is empty, so we can't build redirect url")
	errEmptyCredentials = errors.New("spotify application credentials are empty")
)
