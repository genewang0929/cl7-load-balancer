FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o lb cmd/lb/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/lb .
EXPOSE 8080
CMD ["./lb"]
