override RELEASE="$(git tag -l --points-at HEAD)"
override COMMIT="$(shell git rev-parse --short HEAD)"
override BUILD_TIME="$(shell date -u '+%Y-%m-%dT%H:%M:%S')"

all:

run:
	go run ./cmd/...

compose:
	docker-compose up -d --build

exec-sources:
	docker exec -it auth.sources bash

image:
	docker build \
        --target auth \
		--build-arg RELEASE=${RELEASE} \
		--build-arg COMMIT=${COMMIT} \
		--build-arg BUILD_TIME=${BUILD_TIME} \
		-t "musicmash/auth:latest" .

ensure-go-migrate-installed:
	bash ./scripts/install-go-migrate.sh

db-generate:
	sqlc generate

# show latest applied migration
db-status: ensure-go-migrate-installed
	migrate -path migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose version

# apply migration up
db-up: ensure-go-migrate-installed
	migrate -path migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose up

# apply migration down
db-down: ensure-go-migrate-installed
	migrate -path migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose down
