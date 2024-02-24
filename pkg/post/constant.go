package post

import (
	"context"

	b "github.com/forumGamers/octo-cats/pkg/base"
	protobuf "github.com/forumGamers/octo-cats/protobuf/post"
	tp "github.com/forumGamers/octo-cats/third-party"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepo interface {
	Create(ctx context.Context, data *Post) error
	FindById(ctx context.Context, id primitive.ObjectID, data *Post) error
	GetSession() (mongo.Session, error)
	DeleteOne(ctx context.Context, id primitive.ObjectID) error
	CreateMany(ctx context.Context, datas []any) (*mongo.InsertManyResult, error)
}

type PostRepoImpl struct{ b.BaseRepo }

type PostService interface {
	InsertManyAndBindIds(ctx context.Context, datas []Post) error
	GetPostTags(text string) []string
	CreatePostPayload(userId, text, privacy string, allowComment bool, media []Media, tags []string) Post
	UploadPostMedia(ctx context.Context, file *protobuf.FileHeader) (Media, error)
}

type PostServiceImpl struct {
	Repo PostRepo
	Ik   tp.ImagekitService
}
