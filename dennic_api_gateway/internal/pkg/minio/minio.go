package minio

import (
	"bytes"
	"context"
	"dennic_api_gateway/internal/pkg/config"
	"fmt"
	"net/http"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func UploadToMinio(cfg *config.Config, objectName string, content []byte, bucketName string) (string, error) {
	minioClient, err := minio.New(cfg.MinioService.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioService.AccessKey, cfg.MinioService.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return "", err
	}

	found, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		return "", err
	}
	if !found {
		err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return "", err
		}
		fmt.Println("Bucket created successfully.")
	} else {
		fmt.Println("Bucket already exists.")
	}

	contentType := http.DetectContentType(content)

	opts := minio.PutObjectOptions{ContentType: contentType, UserMetadata: map[string]string{"x-amz-acl": "public-read"}}

	_, err = minioClient.PutObject(context.Background(), bucketName, objectName, bytes.NewReader(content), int64(len(content)), opts)
	if err != nil {
		return "", err
	}

	objectURL := fmt.Sprintf("%s/%s/%s", cfg.MinioService.ImageURL, bucketName, objectName)

	return objectURL, nil
}
