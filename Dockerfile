FROM golang:1.12-alpine

ADD . /app

WORKDIR /app

RUN apk add bash git libc-dev

RUN CGO_ENABLED=0 GOOS=linux

ARG GO111MODULE=on

RUN go build .

ENTRYPOINT ./main