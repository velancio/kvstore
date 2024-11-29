package main

import (
	inmemorystore "censys/internal/kvstore/inmemory"
	"censys/pkg/transport"
	kvstore "censys/proto/gen/proto"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"sync"
)

func main() {
	// Load config from .env
	grpcConnection := fmt.Sprintf(":%s", os.Getenv("KVSTORE_PORT"))
	listen, err := net.Listen("tcp", grpcConnection)
	if err != nil {
		log.Fatalf("Failed to listen: %s", err)
	}

	// Create a gRPC server
	serverRegistrar := grpc.NewServer()
	service := &transport.KvStoreServer{
		Store: &inmemorystore.InMemoryStore{
			Data: sync.Map{},
		},
	}

	// Register the gRPC server
	kvstore.RegisterKvStoreServiceServer(serverRegistrar, service)

	// Start the gRPC server
	if err = serverRegistrar.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}
