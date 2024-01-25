package main

import (
	"log"
	"net"
	"os"

	cc "github.com/forumGamers/octo-cats/controllers"
	"github.com/forumGamers/octo-cats/errors"
	"github.com/forumGamers/octo-cats/protobuf"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	errors.PanicIfError(godotenv.Load())

	address := os.Getenv("PORT")
	if address == "" {
		address = "50052"
	}

	lis, err := net.Listen("tcp", ":"+address)
	if err != nil {
		log.Fatalf("Failed to listen : %s", err.Error())
	}

	grpcServer := grpc.NewServer()
	protobuf.RegisterPostServiceServer(grpcServer, &cc.PostService{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to server : %s", err.Error())
	} else {
		log.Printf("Server is listening")
	}
}
