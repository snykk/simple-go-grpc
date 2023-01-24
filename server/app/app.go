package app

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/snykk/simple-go-grpc/proto/pb"
)

type GolibServer struct {
	pb.GolibServiceServer
}

// unary
func (s *GolibServer) CheckHealty(ctx context.Context, req *pb.Empty) (*pb.Response, error) {
	return &pb.Response{
		Message: "oke...",
	}, nil
}

// bidirectional streaming
func (s *GolibServer) BidirectionalStreaming(stream pb.GolibService_BidirectionalStreamingServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("got book request with title : %v", req.Book.Title)
		res := &pb.BookResponse{
			Book: &pb.Book{
				Title:  req.Book.Title,
				Author: req.Book.Author,
				Desc:   req.Book.Desc,
			},
		}
		if err := stream.Send(res); err != nil {
			return err
		}
	}
}

// client streaming
func (s *GolibServer) ClientStreaming(stream pb.GolibService_ClientStreamingServer) error {
	var books []*pb.Book
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.BooksResponse{Books: books})
		}
		if err != nil {
			return err
		}
		log.Printf("got book request with title : %v", req.Book.Title)
		books = append(books, &pb.Book{
			Title:  req.Book.Title,
			Author: req.Book.Author,
			Desc:   req.Book.Desc,
		})
	}
}

// server streaming
func (s *GolibServer) ServerStreaming(req *pb.BooksRequest, stream pb.GolibService_ServerStreamingServer) error {
	log.Printf("got book request with length: %v", len(req.Books))
	for _, book := range req.Books {
		res := &pb.BookResponse{
			Book: &pb.Book{
				Title:  book.Title,
				Author: book.Author,
				Desc:   book.Desc,
			},
		}
		if err := stream.Send(res); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}
