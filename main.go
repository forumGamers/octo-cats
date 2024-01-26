package main

import (
	"log"
	"net"
	"os"

	cfg "github.com/forumGamers/octo-cats/config"
	cc "github.com/forumGamers/octo-cats/controllers"
	"github.com/forumGamers/octo-cats/errors"
	"github.com/forumGamers/octo-cats/interceptors"
	"github.com/forumGamers/octo-cats/pkg/comment"
	"github.com/forumGamers/octo-cats/pkg/like"
	"github.com/forumGamers/octo-cats/pkg/post"
	"github.com/forumGamers/octo-cats/pkg/share"
	"github.com/forumGamers/octo-cats/protobuf"
	tp "github.com/forumGamers/octo-cats/third-party"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	errors.PanicIfError(godotenv.Load())
	cfg.Connection()

	address := os.Getenv("PORT")
	if address == "" {
		address = "50052"
	}

	lis, err := net.Listen("tcp", ":"+address)
	if err != nil {
		log.Fatalf("Failed to listen : %s", err.Error())
	}

	//thirdparty
	ik := tp.NewImageKit()

	//repository
	postRepo := post.NewPostRepo()
	likeRepo := like.NewLikeRepo()
	commentRepo := comment.NewCommentRepo()
	shareRepo := share.NewShareRepo()

	//services
	postService := post.NewPostService(postRepo, ik)

	interceptor := interceptors.NewInterCeptor()
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptor.Logging, interceptor.UnaryAuthentication),
	)
	protobuf.RegisterPostServiceServer(grpcServer, &cc.PostService{
		GetUser:     interceptor.GetUserFromCtx,
		PostRepo:    postRepo,
		PostService: postService,
		Ik:          ik,
		LikeRepo:    likeRepo,
		CommentRepo: commentRepo,
		ShareRepo:   shareRepo,
	})

	log.Printf("Starting to serve in port : %s", address)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve : %s", err.Error())
	}
}
