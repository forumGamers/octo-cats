package post

import (
	"path/filepath"

	"github.com/forumGamers/octo-cats/errors"
	protobuf "github.com/forumGamers/octo-cats/protobuf/post"
	"google.golang.org/grpc/codes"
)

func ValidateFile(file *protobuf.FileHeader) (string, error) {
	imgExt := []string{"png", "jpg", "jpeg", "gif", "bmp"}
	vidExt := []string{"mp4", "avi", "mkv", "mov"}

	ext := filepath.Ext(file.Filename)[1:]

	for _, val := range imgExt {
		if val == ext {
			if file.Size > 10*1024*1024 {
				return "image", errors.NewAppError(codes.InvalidArgument, "file cannot be larger than 10 mb")
			}
			return "image", nil
		}
	}

	for _, val := range vidExt {
		if val == ext {
			if file.Size > 10*1024*1024 {
				return "video", errors.NewAppError(codes.InvalidArgument, "file cannot be larger than 10 mb")
			}
			return "video", nil
		}
	}
	return "", errors.NewAppError(codes.InvalidArgument, "file type is not supported")
}
