package main

import (
	"log"
	"net"

	"github.com/tjroeder/censys_take_home/cache/grpcserver"
	"github.com/tjroeder/censys_take_home/cache/internal/cache"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	cache := cache.NewCache()
	cacheServer := grpcserver.New(cache)

	// Register the gRPC server with the cacheServer
	grpcserver.RegisterCacheServer(s, cacheServer)

	log.Println("gRPC server listening on :50051")
	if err = s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
