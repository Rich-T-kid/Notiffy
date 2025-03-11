FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

COPY . .

# Build gRPC server
RUN go build -o grpc_server ./cmd/grpc_server/main.go

# Build HTTP server
RUN go build -o http_server ./cmd/http_server/main.go

# make sure permisions are correct I.E is executable
RUN chmod +x /app/run_servers.sh

# Expose ports for both servers
EXPOSE 9999 50051

# Run the script to start both servers
CMD ["/app/run_servers.sh"]
