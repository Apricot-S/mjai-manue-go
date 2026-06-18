# syntax=docker/dockerfile:1

ARG GO_VERSION=1.26
ARG OS_VERSION=trixie

FROM golang:${GO_VERSION}-${OS_VERSION} AS builder

WORKDIR /src

ENV GOEXPERIMENT=jsonv2

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,target=. \
    go build -ldflags="-s -w" -trimpath -o /usr/local/bin/mjai-manue ./cmd/mjai-manue

FROM debian:${OS_VERSION}-slim

ARG UID=10001
RUN useradd \
    --uid "${UID}" \
    --shell "/usr/sbin/nologin" \
    --home-dir "/nonexistent" \
    --no-create-home \
    appuser

WORKDIR /app

COPY --from=builder /usr/local/bin/mjai-manue /usr/local/bin/mjai-manue
COPY scripts/self-match/manue/run_manue.sh .
RUN chmod +x /app/run_manue.sh

USER appuser

ENTRYPOINT ["/app/run_manue.sh"]
