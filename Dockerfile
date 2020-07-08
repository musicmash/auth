FROM golang:latest as builder

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

ARG RELEASE=unset
ARG COMMIT=unset
ARG BUILD_TIME=unset
ENV PROJECT=/go/src/github.com/musicmash/auth/internal

WORKDIR /go/src/github.com/musicmash/auth
COPY cmd cmd
COPY internal internal
COPY go.mod go.mod

RUN go build -v -a \
    -installsuffix cgo \
    -gcflags "all=-trimpath=$(GOPATH)" \
    -ldflags '-linkmode external -extldflags "-static" -s -w \
       -X ${PROJECT}/version.Release=${RELEASE} \
       -X ${PROJECT}/version.Commit=${COMMIT} \
       -X ${PROJECT}/version.BuildTime=${BUILD_TIME}"' \
    -o /usr/local/bin/auth ./cmd/...

FROM alpine:latest

RUN addgroup -S auth && adduser -S auth -G auth
USER auth
WORKDIR /home/auth

COPY --from=builder --chown=auth:auth /usr/local/bin/auth /usr/local/bin/auth

ENTRYPOINT ["/usr/local/bin/auth"]
CMD []