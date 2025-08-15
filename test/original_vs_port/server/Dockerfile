# syntax=docker/dockerfile:1

ARG RUBY_VERSION=3.4.5
ARG OS_VERSION=trixie

FROM ruby:${RUBY_VERSION}-${OS_VERSION}

ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser

RUN gem install mjai

WORKDIR /app

USER appuser

COPY ./run_server.sh .

ENTRYPOINT [ "/app/run_server.sh" ]
