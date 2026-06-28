# syntax=docker/dockerfile:1

ARG OS_VERSION=trixie

FROM gcc:16-${OS_VERSION} AS builder

ARG UID=10001
RUN useradd \
    --uid "${UID}" \
    --shell "/usr/sbin/nologin" \
    --home-dir "/nonexistent" \
    --no-create-home \
    appuser

RUN rm -f /etc/apt/apt.conf.d/docker-clean; \
    echo 'Binary::apt::APT::Keep-Downloaded-Packages "true";' > /etc/apt/apt.conf.d/keep-cache
RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
        libboost-system1.88-dev

RUN mkdir -p /app && \
    chown -R appuser /app

WORKDIR /app

USER appuser

RUN mkdir build
RUN git clone --depth=1 https://github.com/critter-mj/akochan.git

RUN --mount=type=bind,source=patch.sh,target=patch.sh \
    cd ./akochan && \
    /app/patch.sh && \
    cd ./ai_src && \
    make -f Makefile_Linux && \
    cd ../ && \
    make -f Makefile_Linux && \
    cp libai.so system.exe setup_match.json setup_mjai.json ../build && \
    cp -r params ../build

FROM debian:${OS_VERSION}-slim AS final

ARG UID=10001
RUN useradd \
    --uid "${UID}" \
    --shell "/usr/sbin/nologin" \
    --home-dir "/nonexistent" \
    --no-create-home \
    appuser

RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \
    rm -f /etc/apt/apt.conf.d/docker-clean; \
    echo 'Binary::apt::APT::Keep-Downloaded-Packages "true";' > /etc/apt/apt.conf.d/keep-cache && \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
        libboost-system1.88.0 \
        libgomp1 && \
    mkdir -p /app && \
    chown -R appuser /app

WORKDIR /app/akochan

USER appuser

COPY --from=builder /app/build /app/akochan

ENV LD_LIBRARY_PATH=${LD_LIBRARY_PATH:+$LD_LIBRARY_PATH:}/app/akochan

COPY ./run_akochan.sh .

ENTRYPOINT [ "/app/akochan/run_akochan.sh" ]
