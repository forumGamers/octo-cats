package reply

import "github.com/forumGamers/octo-cats/pkg/comment"

type ReplyService interface {
	CreatePayload(text, userId string) comment.ReplyComment
}

type ReplyServiceImpl struct{ Repo comment.CommentRepo }
