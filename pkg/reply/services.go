package reply

import (
	"time"

	"github.com/forumGamers/octo-cats/pkg/comment"
)

func NewReplyService(repo comment.CommentRepo) ReplyService {
	return &ReplyServiceImpl{repo}
}

func (rs *ReplyServiceImpl) CreatePayload(text, userId string) comment.ReplyComment {
	return comment.ReplyComment{
		UserId:    userId,
		Text:      text,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
