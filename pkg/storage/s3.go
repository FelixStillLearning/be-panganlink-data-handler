package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Helper struct {
	client *s3.Client
	bucket string
}

// NewS3Helper inisialisasi helper S3 (belum diaktifkan di router/handler utama)
// Akan membaca kredensial secara otomatis dari environment variables:
// AWS_ACCESS_KEY_ID & AWS_SECRET_ACCESS_KEY
func NewS3Helper(ctx context.Context, region, bucket string) (*S3Helper, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("gagal meload konfigurasi aws: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	return &S3Helper{
		client: client,
		bucket: bucket,
	}, nil
}

// UploadFile mengunggah file gambar/dokumen ke bucket S3
func (h *S3Helper) UploadFile(ctx context.Context, file io.Reader, fileName string, contentType string) (string, error) {
	_, err := h.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(h.bucket),
		Key:         aws.String(fileName),
		Body:        file,
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return "", fmt.Errorf("gagal mengunggah file ke s3: %w", err)
	}

	// Ini adalah format URL S3 standard. Bisa disesuaikan jika nanti pakai CloudFront
	fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", h.bucket, fileName)
	return fileURL, nil
}
