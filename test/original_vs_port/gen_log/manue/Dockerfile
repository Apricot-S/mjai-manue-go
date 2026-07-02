# syntax=docker/dockerfile:1

ARG NODE_VERSION=24
ARG OS_VERSION=trixie

FROM node:${NODE_VERSION}-${OS_VERSION} AS downloader

RUN mkdir -p /app && \
    chown -R node /app

WORKDIR /app

USER node

RUN git clone --depth=1 https://github.com/gimite/mjai-manue.git

ENV NODE_ENV=production

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

RUN npm install -g coffeescript@1.12.7

COPY --from=downloader /app /app

# Patch only getHoraFactorsDist(actor). getScoreChangesDistOnHora(metric) has
# the same expression, but there @player() is the winner and the expression is
# already correct.
RUN sed -i '/getHoraFactorsDist: (actor) ->/,/return new ProbDist(m)/ s/target == @player()/target == actor/' \
    /app/mjai-manue/coffee/manue_ai.coffee

USER node

COPY ./run_manue.sh .

ENTRYPOINT [ "/app/run_manue.sh" ]
