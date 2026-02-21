# syntax=docker/dockerfile:1

ARG RUBY_VERSION=4.0.1
ARG OS_VERSION=trixie

FROM ruby:${RUBY_VERSION}-${OS_VERSION} AS builder

WORKDIR /opt

RUN git clone --depth=1 https://github.com/gimite/mjai.git

RUN --mount=type=bind,source=patch.sh,target=patch.sh \
    cd mjai && \
    /opt/patch.sh && \
    gem build mjai.gemspec && \
    gem install mjai-*.gem

FROM ruby:${RUBY_VERSION}-slim-${OS_VERSION} AS final

ARG UID=10001
RUN useradd \
    --uid "${UID}" \
    --shell "/usr/sbin/nologin" \
    --home-dir "/nonexistent" \
    --no-create-home \
    appuser

COPY --from=builder /usr/local/bundle /usr/local/bundle
RUN chown -R appuser:appuser /usr/local/bundle/gems/mjai-0.0.7/share/html/css

WORKDIR /app

USER appuser

COPY ./run_server.sh .

ENTRYPOINT [ "/app/run_server.sh" ]
