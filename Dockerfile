FROM golang:1.22.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /mess cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /mess .
COPY .env .
COPY wait-for-kafka.sh .
RUN chmod +x wait-for-kafka.sh

CMD ["./wait-for-kafka.sh", "kafka", "./mess"]
