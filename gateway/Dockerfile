#builder
FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY gateway/ ./gateway/
COPY pkg/ ./pkg/
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway-bin ./gateway/cmd/main.go

#runtime
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/gateway-bin .

EXPOSE 8080

CMD ["./gateway-bin"]