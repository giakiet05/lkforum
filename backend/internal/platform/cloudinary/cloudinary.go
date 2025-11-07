package cloudinary

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/giakiet05/lkforum/internal/config"
)

func newCld() (*cloudinary.Cloudinary, error) {
	return cloudinary.NewFromParams(config.Cfg.Cloudinary.CloudName, config.Cfg.Cloudinary.APIKey, config.Cfg.Cloudinary.APISecret)
}

func Upload(file multipart.File, fileHeader *multipart.FileHeader) (*uploader.UploadResult, error) {
	cld, err := newCld()
	if err != nil {
		return nil, err
	}

	return cld.Upload.Upload(context.Background(), file, uploader.UploadParams{
		Folder: config.Cfg.Cloudinary.UploadFolder,
	})
}

func Delete(publicID string) (*uploader.DestroyResult, error) {
	cld, err := newCld()
	if err != nil {
		return nil, err
	}

	return cld.Upload.Destroy(context.Background(), uploader.DestroyParams{PublicID: publicID})
}
