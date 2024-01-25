package thirdparty

import (
	"context"
	"os"

	"github.com/codedius/imagekit-go"
	"github.com/forumGamers/octo-cats/errors"
)

func NewImageKit() ImagekitService {
	opts := imagekit.Options{
		PublicKey:  os.Getenv("IMAGEKIT_PUBLIC_KEY"),
		PrivateKey: os.Getenv("IMAGEKIT_PRIVATE_KEY"),
	}

	ik, err := imagekit.NewClient(&opts)
	errors.PanicIfError(err)

	return &ImagekitServiceImpl{ik}
}

func (ik *ImagekitServiceImpl) UploadFile(ctx context.Context, upload UploadFile) (*imagekit.UploadResponse, error) {
	return ik.Client.Upload.ServerUpload(ctx, &imagekit.UploadRequest{
		File:              upload.Data,
		FileName:          upload.Name,
		UseUniqueFileName: true,
		Folder:            upload.Folder,
	})
}

func (ik *ImagekitServiceImpl) DeleteFile(ctx context.Context, fileId string) error {
	return ik.Client.Media.DeleteFile(ctx, fileId)
}

func (ik *ImagekitServiceImpl) DeleteBulkFile(ctx context.Context, fileIds []string) error {
	_, err := ik.Client.Media.DeleteFiles(ctx, &imagekit.DeleteFilesRequest{FileIDs: fileIds})
	return err
}
