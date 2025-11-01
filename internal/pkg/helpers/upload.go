package helpers

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

type CloudinaryUploader struct {
	cld *cloudinary.Cloudinary
}

func New(cld *cloudinary.Cloudinary) *CloudinaryUploader {
	return &CloudinaryUploader{cld: cld}
}

type UploadOptions struct {
	Folder string
}

type UploadedImage struct {
	URL      string
	PublicID string
	FileName string
	FileType string
	FileSize int64
	Width    int
	Height   int
}

func (cu *CloudinaryUploader) UploadImage(file multipart.File, fileHeader *multipart.FileHeader, opts *UploadOptions) (*UploadedImage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ext := filepath.Ext(fileHeader.Filename)
	publicID := strings.TrimSuffix(uuid.New().String(), ext)

	result, err := cu.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID:       publicID,
		Folder:         opts.Folder,
		ResourceType:   "image",
		Transformation: "q_auto,f_auto",
	})
	if err != nil {
		return nil, fmt.Errorf("Cloudinary upload failed: %w", err)
	}

	return &UploadedImage{
		URL:      result.SecureURL,
		PublicID: result.PublicID,
		FileName: fileHeader.Filename,
		FileType: fileHeader.Header.Get("Content-Type"),
		FileSize: fileHeader.Size,
		Width:    result.Width,
		Height:   result.Height,
	}, nil
}

func (cu *CloudinaryUploader) Delete(publicID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := cu.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: "image",
	})
	return err
}
