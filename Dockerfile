FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . . 
RUN go build -v -o /app/main ./cmd


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main /app/main
COPY --from=builder /app/config.yaml /app/config.yaml

EXPOSE 8000
CMD ["/app/main"]