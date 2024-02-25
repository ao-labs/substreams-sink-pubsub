FROM golang:1.22-alpine as build
WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o /app/substreams-sink-pubsub ./cmd/substreams-sink-pubsub/*

