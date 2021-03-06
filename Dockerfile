FROM golang:1-alpine as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

ARG RELEASE=unset
ARG COMMIT=unset
ARG BUILD_TIME=unset
ENV PROJECT=github.com/musicmash/auth

WORKDIR /go/src/github.com/musicmash/auth
COPY migrations /var/auth/migrations
COPY cmd cmd
COPY internal internal
COPY go.mod go.mod
COPY go.sum go.sum

RUN go build -v -a \
    -gcflags "all=-trimpath=${WORKDIR}" \
    -ldflags "-w -s \
       -X ${PROJECT}/internal/version.Release=${RELEASE} \
       -X ${PROJECT}/internal/version.Commit=${COMMIT} \
       -X ${PROJECT}/internal/version.BuildTime=${BUILD_TIME}" \
    -o /usr/local/bin/auth ./cmd/...

FROM alpine:latest as auth

RUN addgroup -S auth && adduser -S auth -G auth
USER auth
WORKDIR /home/auth

COPY --from=builder --chown=auth:auth /var/auth/migrations /var/auth/migrations
COPY --from=builder --chown=auth:auth /usr/local/bin/auth /usr/local/bin/auth

ENTRYPOINT ["/usr/local/bin/auth"]
CMD []