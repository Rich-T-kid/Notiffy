#!/bin/sh

# Start gRPC server in the background
./grpc_server &

# Start HTTP server in the background
./http_server &

# Wait for all background processes to finish
wait

