package minio

import (
	"Healthcare_Evrone/internal/pkg/config"
	"strings"
)

var cfg = config.New()

func AddImageUrl(imageUrl, bucketName string) string {
	str := cfg.MinioService.Endpoint + "/" + bucketName + "/" + imageUrl
	return str
}

func RemoveImageUrl(imageUrl string) string {
	str := strings.Split(imageUrl, "/")
	return str[len(str)-1]
}
