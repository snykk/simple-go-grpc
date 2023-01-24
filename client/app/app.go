package app

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/snykk/simple-go-grpc/proto/pb"
)

type GolibClient struct {
	Client pb.GolibServiceClient
}

// unary
func (a *GolibClient) CallCheckHealty() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := a.Client.CheckHealty(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("somethig when wrong: %v", err)
	}
	log.Printf("%s", res.Message)
}

// bidireectional streaming
func (a *GolibClient) CallBidirectionalStream(books *pb.BooksRequest) {
	log.Printf("starting bidirectional streaming")
	stream, err := a.Client.BidirectionalStreaming(context.Background())
	if err != nil {
		log.Fatalf("could not send books: %v", err)
	}

	go func() {
		for {
			message, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while streaming %v", err)
			}
			fmt.Println("\nBook Data Stream")
			fmt.Printf("Tittle\t\t:%s\n", message.Book.Title)
			fmt.Printf("Author\t\t:%s\n", message.Book.Author)
			fmt.Printf("Description\t:%s\n", message.Book.Desc)
		}
	}()

	for _, book := range books.Books {
		req := &pb.BookRequest{
			Book: &pb.Book{
				Title:  book.Title,
				Author: book.Author,
				Desc:   book.Desc,
			},
		}
		if err := stream.Send(req); err != nil {
			log.Fatalf("error while sending %v", err)
		}
		time.Sleep(2 * time.Second)
	}

	stream.CloseSend()
	log.Printf("bidirectional Streaming finished")
}

// Client streaming
func (a *GolibClient) CallClientStream(books *pb.BooksRequest) {
	log.Printf("starting client streaming")
	stream, err := a.Client.ClientStreaming(context.Background())
	if err != nil {
		log.Fatalf("could not send books: %v", err)
	}

	for _, book := range books.Books {
		req := &pb.BookRequest{
			Book: &pb.Book{
				Title:  book.Title,
				Author: book.Author,
				Desc:   book.Desc,
			},
		}
		if err := stream.Send(req); err != nil {
			log.Fatalf("error while sending %v", err)
		}
		log.Printf("sent book request with title: %s", book.Title)
		time.Sleep(2 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	log.Printf("client Streaming finished")
	if err != nil {
		log.Fatalf("error while receiving stream data %v", err)
	}
	// log.Printf("%v", res.Books)
	for idx, book := range res.Books {
		fmt.Printf("\nBook %d\n", idx+1)
		fmt.Printf("Title\t\t: %s\n", book.Title)
		fmt.Printf("Author\t\t: %s\n", book.Author)
		fmt.Printf("Description\t: %s\n", book.Desc)
	}
}

// server streaming
func (a *GolibClient) CallServerStream(books *pb.BooksRequest) {
	log.Printf("starting server streaming")
	stream, err := a.Client.ServerStreaming(context.Background(), books)
	if err != nil {
		log.Fatalf("could not send books: %v", err)
	}

	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while streaming %v", err)
		}
		fmt.Println("\nBook Data")
		fmt.Printf("Tittle\t\t:%s\n", message.Book.Title)
		fmt.Printf("Author\t\t:%s\n", message.Book.Author)
		fmt.Printf("Description\t:%s\n", message.Book.Desc)
	}

	log.Printf("streaming finished")
}
