package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type AzureHelper struct {
	client        *azblob.Client
	containerName string
}

// NewAzureHelper inisialisasi helper Azure Blob Storage (belum diaktifkan di router/handler utama)
// Gunakan connectionString yang didapatkan dari portal Azure Storage Account.
func NewAzureHelper(connectionString, containerName string) (*AzureHelper, error) {
	client, err := azblob.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, fmt.Errorf("gagal meload konfigurasi azure: %w", err)
	}

	return &AzureHelper{
		client:        client,
		containerName: containerName,
	}, nil
}

// UploadFile mengunggah file gambar/dokumen ke Azure Blob Storage
func (h *AzureHelper) UploadFile(ctx context.Context, file io.Reader, fileName string, contentType string) (string, error) {
	// Menentukan metadata dan content type agar file bisa dibuka di browser
	options := &azblob.UploadStreamOptions{
		HTTPHeaders: &azblob.BlobHTTPHeaders{
			BlobContentType: &contentType,
		},
	}

	_, err := h.client.UploadStream(ctx, h.containerName, fileName, file, options)
	if err != nil {
		return "", fmt.Errorf("gagal mengunggah file ke azure blob: %w", err)
	}

	// Format URL indikatif Azure Blob
	// URL aslinya bisa berbeda jika menggunakan Custom Domain / CDN
	fileURL := fmt.Sprintf("%s%s/%s", h.client.URL(), h.containerName, fileName)
	return fileURL, nil
}
