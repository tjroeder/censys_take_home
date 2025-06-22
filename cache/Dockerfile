# cache/Dockerfile
FROM golang:1.24.3-alpine

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the gRPC server
RUN CGO_ENABLED=0 go build -o grpc-server ./cache/main.go

# Use the binary as the container entrypoint
ENTRYPOINT ["./grpc-server"]