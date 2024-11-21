FROM golang:1.22-alpine AS build
WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o /app/substreams-sink-pubsub ./cmd/substreams-sink-pubsub/*

FROM scratch
WORKDIR /app
COPY --from=build /app/substreams-sink-pubsub ./substreams-sink-pubsub

ENTRYPOINT ["./substreams-sink-pubsub"]
