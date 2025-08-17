# syntax=docker/dockerfile:1

ARG NODE_VERSION=22
ARG OS_VERSION=trixie

FROM node:${NODE_VERSION}-${OS_VERSION} AS downloader

RUN mkdir -p /app && \
    chown -R node /app

WORKDIR /app

USER node

RUN git clone https://github.com/gimite/mjai-manue.git

RUN cd mjai-manue/coffee && \
    npm install \
        glob@11.0.3 \
        printf@0.6.1 \
        seed-random@2.2.0

FROM node:${NODE_VERSION}-${OS_VERSION}-slim AS final

ENV NODE_ENV=production

RUN mkdir -p /app && \
    chown -R node /app

WORKDIR /app

COPY --from=downloader /app /app

RUN npm install -g coffeescript@1.12.7

USER node

COPY ./run_manue.sh .

ENTRYPOINT [ "/app/run_manue.sh" ]
