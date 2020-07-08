COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%dT%H:%M:%S')

all:

run:
	go run ./cmd/...

image:
	docker build \
		--build-arg COMMIT=${COMMIT} \
		--build-arg BUILD_TIME="$(BUILD_TIME)" \
		--compress \
		-t auth:latest .
