# api/Dockerfile
FROM golang:1.24.3-alpine

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the api gateway
RUN CGO_ENABLED=0 go build -o api-gateway ./api/main.go

ENTRYPOINT ["./api-gateway"]