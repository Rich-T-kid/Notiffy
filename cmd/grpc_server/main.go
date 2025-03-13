package main

import (
	// this package needs to always to be run first b4 all other custom packages
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	grpcService "github.com/Rich-T-kid/Notiffy/api/grpc"
	"github.com/Rich-T-kid/Notiffy/api/grpc/protobuff"
	_ "github.com/Rich-T-kid/Notiffy/internal/enviroment" // this package needs to always to be run first b4 all other custom packages
)

// TODO: for now just implement the GRPC server first. We can add and HTTP one later. when http server is implemented seperate main function into two serpate main.go files within a cmd subdirectory.
var port = os.Getenv("GRPC_SERVER_PORT")

func main() {
	// when i finish the server implementation this is where it will go

	serverImplementation := grpcService.NewGServer()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcService.UnaryInterceptor))

	protobuff.RegisterNotificationServiceServer(grpcServer, serverImplementation)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to create listener: %v", err)
	}

	fmt.Printf("Starting gRPC server on port %s...\n", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}

}
