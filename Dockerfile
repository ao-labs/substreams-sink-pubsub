FROM golang:1.22-alpine AS build
WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/substreams-sink-pubsub ./cmd/substreams-sink-pubsub

FROM scratch
WORKDIR /app

# Import the certificate authority certificates.
# Since scratch images do not contain any CA certificates, we need to import these from the builder image.
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /app/substreams-sink-pubsub ./substreams-sink-pubsub

ENTRYPOINT ["./substreams-sink-pubsub"]
