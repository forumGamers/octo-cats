package post

import (
	"context"
	"fmt"
	"strings"
	"time"

	protobuf "github.com/forumGamers/octo-cats/protobuf/post"
	tp "github.com/forumGamers/octo-cats/third-party"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewPostService(repo PostRepo, ik tp.ImagekitService) PostService {
	return &PostServiceImpl{repo, ik}
}

func (s *PostServiceImpl) InsertManyAndBindIds(ctx context.Context, datas []Post) error {
	var payload []any

	for _, data := range datas {
		payload = append(payload, data)
	}

	ids, err := s.Repo.CreateMany(ctx, payload)
	if err != nil {
		return err
	}

	for i := 0; i < len(ids.InsertedIDs); i++ {
		id := ids.InsertedIDs[i].(primitive.ObjectID)
		datas[i].Id = id
	}
	return nil
}

func (s *PostServiceImpl) GetPostTags(text string) []string {
	modified := text
	for _, p := range "!@#$%^&*)(_=+?.,;:'" {
		modified = strings.ReplaceAll(modified, string(p), " ")
	}
	return strings.Split(modified, " ")
}

func (s *PostServiceImpl) CreatePostPayload(userId, text, privacy string, allowComment bool, media []Media, tags []string) Post {
	return Post{
		UserId:       userId,
		Text:         text,
		Media:        media,
		AllowComment: allowComment,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Tags:         tags,
		Privacy:      privacy,
	}
}

func (s *PostServiceImpl) UploadPostMedia(ctx context.Context, file *protobuf.FileHeader) (data Media, err error) {
	fileType, err := ValidateFile(file)
	if err != nil {
		return
	}

	response, err := s.Ik.UploadFile(ctx, tp.UploadFile{
		Data: file.Content, Name: file.Filename, Folder: fmt.Sprintf("post_%s", fileType),
	})
	if err != nil {
		return
	}

	data = Media{
		Url:  response.URL,
		Id:   response.FileID,
		Type: fileType,
	}
	return
}
