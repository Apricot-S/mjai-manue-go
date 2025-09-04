# syntax=docker/dockerfile:1

ARG GO_VERSION=1.25
ARG OS_VERSION=trixie

FROM golang:${GO_VERSION}-${OS_VERSION} AS builder

WORKDIR /app

ENV GOEXPERIMENT=jsonv2

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,target=. \
    go build -ldflags="-s -w" -trimpath -o /usr/local/bin/ ./cmd/mjai-manue

FROM debian:${OS_VERSION}-slim

RUN rm -f /etc/apt/apt.conf.d/docker-clean; \
    echo 'Binary::apt::APT::Keep-Downloaded-Packages "true";' > /etc/apt/apt.conf.d/keep-cache
RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
        zip

WORKDIR /build

COPY /scripts/mjai.app/bot.py .

COPY --from=builder /usr/local/bin/mjai-manue .

RUN zip mjai-app.zip *
