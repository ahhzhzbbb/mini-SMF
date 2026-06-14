#builder
FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY pdu-session/ ./pdu-session/
COPY pkg/ ./pkg/
RUN CGO_ENABLED=0 GOOS=linux go build -o pdu-session-bin ./pdu-session/cmd/main.go

#runtime
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/pdu-session-bin .

EXPOSE 8081

CMD ["./pdu-session-bin"]