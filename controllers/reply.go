package controllers

import (
	"context"

	"github.com/forumGamers/octo-cats/pkg/comment"
	"github.com/forumGamers/octo-cats/pkg/reply"
	"github.com/forumGamers/octo-cats/pkg/user"
	protobuf "github.com/forumGamers/octo-cats/protobuf/reply"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ReplyService struct {
	protobuf.UnimplementedReplyServiceServer
	GetUser        func(ctx context.Context) user.User
	CommentRepo    comment.CommentRepo
	CommentService comment.CommentService
	ReplyService   reply.ReplyService
}

func (s *ReplyService) CreateReply(ctx context.Context, req *protobuf.CommentForm) (*protobuf.Reply, error) {
	if req.Text == "" {
		return nil, status.Error(codes.InvalidArgument, "text is required")
	}

	commentId, err := primitive.ObjectIDFromHex(req.CommentId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid CommentId")
	}

	var commentData comment.Comment
	if err := s.CommentRepo.FindById(ctx, commentId, &commentData); err != nil {
		return nil, err
	}

	replyPayload := s.ReplyService.CreatePayload(req.Text, s.GetUser(ctx).UUID)
	if err := s.CommentRepo.CreateReply(ctx, commentId, &replyPayload); err != nil {
		return nil, err
	}

	return &protobuf.Reply{
		XId:       replyPayload.Id.Hex(),
		Text:      replyPayload.Text,
		UserId:    replyPayload.UserId,
		CreatedAt: replyPayload.CreatedAt.Local().String(),
		UpdatedAt: replyPayload.UpdatedAt.Local().String(),
	}, nil
}
