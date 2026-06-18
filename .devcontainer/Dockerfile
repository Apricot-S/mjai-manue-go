FROM mcr.microsoft.com/devcontainers/go:2-1.26-trixie

COPY --from=ghcr.io/astral-sh/uv:0.11.21 /uv /uvx /bin/

RUN apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
        gnuplot \
        fonts-ipafont-gothic \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

ENV UV_LINK_MODE=copy
