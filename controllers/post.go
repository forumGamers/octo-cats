package controllers

import (
	"context"
	"sync"

	"github.com/forumGamers/octo-cats/pkg/comment"
	"github.com/forumGamers/octo-cats/pkg/like"
	"github.com/forumGamers/octo-cats/pkg/post"
	"github.com/forumGamers/octo-cats/pkg/share"
	"github.com/forumGamers/octo-cats/pkg/user"
	protobuf "github.com/forumGamers/octo-cats/protobuf/post"
	tp "github.com/forumGamers/octo-cats/third-party"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostService struct {
	protobuf.UnimplementedPostServiceServer
	GetUser     func(ctx context.Context) user.User
	PostRepo    post.PostRepo
	PostService post.PostService
	Ik          tp.ImagekitService
	LikeRepo    like.LikeRepo
	CommentRepo comment.CommentRepo
	ShareRepo   share.ShareRepo
}

func (s *PostService) CreatePost(ctx context.Context, req *protobuf.PostForm) (*protobuf.Messages, error) {
	println("masuk")
	return nil, nil
}

func (s *PostService) DeletePost(ctx context.Context, req *protobuf.PostIdPayload) (*protobuf.Messages, error) {
	if req.XId == "" {
		return nil, status.Error(codes.InvalidArgument, "_id is required")
	}

	postId, err := primitive.ObjectIDFromHex(req.XId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid ObjectId")
	}

	var data post.Post
	if err := s.PostRepo.FindById(ctx, postId, &data); err != nil {
		return nil, err
	}

	user := s.GetUser(ctx)
	if user.UUID != data.UserId && user.LoggedAs != "Admin" {
		return nil, status.Error(codes.Unauthenticated, "Forbidden")
	}

	session, err := s.PostRepo.GetSession()
	if err != nil {
		return nil, status.Error(codes.Unavailable, "Failed get session")
	}
	defer session.EndSession(ctx)

	dbCtx := mongo.NewSessionContext(ctx, session)
	if err := session.StartTransaction(); err != nil {
		return nil, status.Error(codes.Unavailable, "Failed start DB Operations")
	}

	var wg sync.WaitGroup
	errCh := make(chan error)
	handlers := []func(){
		func() {
			defer wg.Done()
			var ids []string

			for _, media := range data.Media {
				ids = append(ids, media.Id)
			}

			if len(ids) > 0 {
				errCh <- s.Ik.DeleteBulkFile(dbCtx, ids)
			} else {
				errCh <- nil
			}
		},
		func() {
			defer wg.Done()
			errCh <- s.LikeRepo.DeletePostLikes(dbCtx, data.Id)
		},
		func() {
			defer wg.Done()
			errCh <- s.PostRepo.DeleteOne(dbCtx, data.Id)
		},
		func() {
			defer wg.Done()
			errCh <- s.ShareRepo.DeleteMany(ctx, data.Id)
		},
		func() {
			defer wg.Done()
			errCh <- s.CommentRepo.DeleteMany(ctx, data.Id)
		},
	}

	for _, handler := range handlers {
		wg.Add(1)
		go handler()
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			session.AbortTransaction(dbCtx)
			return nil, err
		}
	}

	if err := session.CommitTransaction(dbCtx); err != nil {
		session.AbortTransaction(dbCtx)
		return nil, err
	}

	return &protobuf.Messages{Message: "success"}, nil
}
