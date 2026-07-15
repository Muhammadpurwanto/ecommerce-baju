package storage

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageService interface {
	UploadAvatar(ctx context.Context, file *multipart.FileHeader, objectName string) (string, error)
	InitBucket(ctx context.Context) error
}

type minioStorage struct {
	client    *minio.Client
	bucket    string
	publicURL string // Ditambahkan public URL config
}

func NewMinioStorage(cfg *config.Config) (StorageService, error) {
	minioClient, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: false, // ubah ke true jika pakai HTTPS
	})
	if err != nil {
		return nil, err
	}

	return &minioStorage{
		client:    minioClient,
		bucket:    cfg.MinioBucket,
		publicURL: cfg.MinioPublicURL,
	}, nil
}

func (s *minioStorage) InitBucket(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		return err
	}
	if !exists {
		err = s.client.MakeBucket(ctx, s.bucket, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
		
		// Set public policy for avatars
		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Action": ["s3:GetObject"],
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Resource": ["arn:aws:s3:::%s/*"]
				}
			]
		}`, s.bucket)
		
		err = s.client.SetBucketPolicy(ctx, s.bucket, policy)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *minioStorage) UploadAvatar(ctx context.Context, file *multipart.FileHeader, objectName string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = s.client.PutObject(ctx, s.bucket, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	// Menggunakan Public URL dari .env, bukan hardcoded localhost:9000
	return fmt.Sprintf("%s/%s/%s", s.publicURL, s.bucket, objectName), nil
}
