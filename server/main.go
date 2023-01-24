package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/snykk/simple-go-grpc/proto/pb"
	app "github.com/snykk/simple-go-grpc/server/app"
	"google.golang.org/grpc"
)

// define the port
const (
	port = 8080
)

func main() {
	//listen on the port
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to start server %v", err)
	}
	// create a new gRPC server
	grpcServer := grpc.NewServer()
	// register the Golib service
	pb.RegisterGolibServiceServer(grpcServer, &app.GolibServer{})
	log.Printf("server started at %v", lis.Addr())
	//list is the port, the grpc server needs to start there
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to start: %v", err)
	}
}
