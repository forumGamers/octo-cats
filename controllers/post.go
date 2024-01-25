package controllers

import (
	"context"

	"github.com/forumGamers/octo-cats/protobuf"
)

type PostService struct {
	protobuf.UnimplementedPostServiceServer
}

func (s *PostService) CreatePost(ctx context.Context, req *protobuf.PostForm) (*protobuf.Messages, error) {
	println("masuk")
	return nil, nil
}

func (s *PostService) DeletePost(ctx context.Context, req *protobuf.PostIdPayload) (*protobuf.Messages, error) {
	

	return &protobuf.Messages{Message: "success"}, nil
}
