FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o api cmd/api/main.go
RUN go build -o poller cmd/poller/main.go
RUN go build -o consumer cmd/consumer/main.go

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/api .
COPY --from=builder /app/poller .
COPY --from=builder /app/consumer .

EXPOSE 8080

CMD ["./api"]