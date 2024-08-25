FROM golang:1.23.0-bullseye

LABEL maintainer="akamiya208"

ENV APP_DIR /app
WORKDIR ${APP_DIR}

# Requirements are installed here.
