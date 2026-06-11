package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryHelper struct {
	cld *cloudinary.Cloudinary
}

// NewCloudinaryHelper inisialisasi helper Cloudinary
func NewCloudinaryHelper(cloudName, apiKey, apiSecret string) (*CloudinaryHelper, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("gagal menginisialisasi cloudinary: %w", err)
	}

	return &CloudinaryHelper{
		cld: cld,
	}, nil
}

// UploadFile mengunggah file media (gambar/produk) ke Cloudinary
func (h *CloudinaryHelper) UploadFile(ctx context.Context, file io.Reader, fileName string) (string, error) {
	resp, err := h.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:   "panganlink",
		PublicID: fileName,
	})
	if err != nil {
		return "", fmt.Errorf("gagal mengunggah file ke cloudinary: %w", err)
	}

	return resp.SecureURL, nil
}
