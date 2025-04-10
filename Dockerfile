FROM golang:1.24-alpine AS builder

WORKDIR /app

ENV GOCACHE=/go/std/cache
RUN go build -v std

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg
COPY static ./static
RUN go build -v -o ./main ./cmd

COPY config.yaml ./

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main /app/main
COPY --from=builder /app/config.yaml /app/config.yaml

EXPOSE 8000
CMD ["/app/main"]