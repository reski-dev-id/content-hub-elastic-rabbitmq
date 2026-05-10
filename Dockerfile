FROM golang:1.25-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o api cmd/api/main.go
RUN go build -o poller cmd/poller/main.go
RUN go build -o consumer cmd/consumer/main.go

EXPOSE 8080