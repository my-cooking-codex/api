# syntax=docker/dockerfile:1.4
ARG GO_VERSION=1.20
ARG ALPINE_VERSION=3.17

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder

    WORKDIR /usr/src/app

    RUN apk add --no-cache vips-dev gcc musl-dev

    COPY go.mod go.sum ./
    RUN go mod download && go mod verify

    COPY . .

    RUN CGO_ENABLED=1 go build -v -o /usr/local/bin/app

FROM alpine:${ALPINE_VERSION}

    RUN apk add --no-cache vips

    COPY --from=builder --link /usr/local/bin/app /usr/local/bin/app

    ENV BIND__HOST=0.0.0.0
    ENV BIND__PORT=8000
    ENV DB__URI=/data/db.sqlite
    ENV DB__TYPE=sqlite
    ENV DATA_PATH=/data

    VOLUME "/data"
    EXPOSE 8000

    CMD ["app"]
