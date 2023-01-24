package main

import (
	"fmt"
	"log"

	app "github.com/snykk/simple-go-grpc/client/app"
	pb "github.com/snykk/simple-go-grpc/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	port = 8080
)

func main() {
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	app := app.GolibClient{
		Client: pb.NewGolibServiceClient(conn),
	}

	books := &pb.BooksRequest{
		Books: []*pb.Book{
			{
				Title:  "Atomic Habits",
				Author: "James Clear",
				Desc:   "Lorem Ipsum Doler Sit Amet",
			},
			{
				Title:  "Selena",
				Author: "Tere Liye",
				Desc:   "Lorem Ipsum Doler Sit Amet",
			},
			{
				Title:  "Mindset",
				Author: "Carrol Dweck",
				Desc:   "Lorem Ipsum Doler Sit Amet",
			},
		},
	}

	fmt.Printf("\n>>> UNARY <<<\n")
	app.CallCheckHealty()
	fmt.Printf("\n>>> SERVER STREAMING <<<\n")
	app.CallServerStream(books)
	fmt.Printf("\n>>> CLIENT STREAMING <<<\n")
	app.CallClientStream(books)
	fmt.Printf("\n>>> BIDIRECTIONAL STREAMING <<<\n")
	app.CallBidirectionalStream(books)
}
