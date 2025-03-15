FROM golang:1.24-bookworm AS builder

WORKDIR /app

COPY . .

RUN go build -ldflags="-s -w" -trimpath ./cmd/mjai-manue

FROM debian:bookworm-slim

RUN rm -f /etc/apt/apt.conf.d/docker-clean; \
    echo 'Binary::apt::APT::Keep-Downloaded-Packages "true";' > /etc/apt/apt.conf.d/keep-cache
RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
        zip

WORKDIR /build

COPY --from=builder /app/mjai-manue .

COPY /scripts/mjai.app/bot.py .

RUN zip mjai-app.zip *
