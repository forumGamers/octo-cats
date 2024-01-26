package post

import tp "github.com/forumGamers/octo-cats/third-party"

func NewPostService(repo PostRepo, ik tp.ImagekitService) PostService {
	return &PostServiceImpl{repo, ik}
}
