FROM docker.io/golang:bookworm

RUN apt update && apt install -y build-essential mpv libmpv-dev
