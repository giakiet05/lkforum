package cloudinary

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
)

func newCld() (*cloudinary.Cloudinary, error) {
	return cloudinary.NewFromParams(config.Cfg.Cloudinary.CloudName, config.Cfg.Cloudinary.APIKey, config.Cfg.Cloudinary.APISecret)
}

func Upload(file multipart.File) (*uploader.UploadResult, error) {
	cld, err := newCld()
	if err != nil {
		return nil, err
	}

	return cld.Upload.Upload(context.Background(), file, uploader.UploadParams{
		Folder: config.Cfg.Cloudinary.UploadFolder,
	})
}

func UploadVideo(file multipart.File) (*uploader.UploadResult, error) {
	cld, err := newCld()
	if err != nil {
		return nil, err
	}

	return cld.Upload.Upload(context.Background(), file, uploader.UploadParams{
		Folder:       config.Cfg.Cloudinary.UploadFolder,
		ResourceType: "video",
	})
}

func Delete(publicID string) (*uploader.DestroyResult, error) {
	cld, err := newCld()
	if err != nil {
		return nil, err
	}

	return cld.Upload.Destroy(context.Background(), uploader.DestroyParams{PublicID: publicID})
}

// UploadImages uploads multiple images and returns model.Image slice.
// This function handles both single and multiple image uploads.
func UploadImages(files []*multipart.FileHeader) ([]*model.Image, error) {
	if len(files) == 0 {
		return nil, errors.New("no images provided")
	}

	var uploadedImages []*model.Image
	var lastErr error

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			lastErr = err
			continue
		}

		result, err := Upload(file)
		file.Close() // Close immediately, not defer in loop

		if err != nil {
			lastErr = err
			continue
		}

		uploadedImages = append(uploadedImages, &model.Image{
			URL:        result.SecureURL,
			PublicID:   result.PublicID,
			UploadedAt: time.Now(),
		})
	}

	if len(uploadedImages) == 0 {
		if lastErr != nil {
			return nil, fmt.Errorf("all image uploads failed: %w", lastErr)
		}
		return nil, errors.New("all image uploads failed")
	}

	return uploadedImages, nil
}

// UploadVideos uploads multiple videos and returns model.Video slice.
func UploadVideos(files []*multipart.FileHeader) ([]*model.Video, error) {
	if len(files) == 0 {
		return nil, errors.New("no videos provided")
	}

	var uploadedVideos []*model.Video
	var lastErr error

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			lastErr = err
			continue
		}

		result, err := UploadVideo(file) // Use UploadVideo
		file.Close()

		if err != nil {
			lastErr = err
			continue
		}

		uploadedVideos = append(uploadedVideos, &model.Video{
			URL:        result.SecureURL,
			PublicID:   result.PublicID,
			UploadedAt: time.Now(),
		})
	}

	if len(uploadedVideos) == 0 {
		if lastErr != nil {
			return nil, fmt.Errorf("all video uploads failed: %w", lastErr)
		}
		return nil, errors.New("all video uploads failed")
	}

	return uploadedVideos, nil
}
