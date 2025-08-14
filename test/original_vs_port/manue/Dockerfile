# syntax=docker/dockerfile:1

ARG NODE_VERSION=22
ARG OS_VERSION=trixie-slim

FROM node:${NODE_VERSION}-${OS_VERSION} AS base

FROM base AS downloader

RUN rm -f /etc/apt/apt.conf.d/docker-clean; \
    echo 'Binary::apt::APT::Keep-Downloaded-Packages "true";' > /etc/apt/apt.conf.d/keep-cache
RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
        ca-certificates \
        git-core

RUN mkdir -p /app && \
    chown -R node /app

WORKDIR /app

USER node

RUN git clone https://github.com/gimite/mjai-manue.git

FROM base AS final

ENV NODE_ENV=production

RUN mkdir -p /app && \
    chown -R node /app

WORKDIR /app

COPY --from=downloader /app /app

RUN cd mjai-manue/coffee && \
    npm install -g coffeescript@1.12.7 && \
    npm install \
        glob@11.0.3 \
        printf@0.6.1 \
        seed-random@2.2.0

USER node

COPY --chown=node ./run_manue.sh .

ENTRYPOINT [ "./run_manue.sh" ]
