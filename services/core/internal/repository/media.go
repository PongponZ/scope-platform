package repository

import (
	"github.com/PongponZ/scope-platform/core/internal/entity"
	"github.com/go-redis/redis/v8"
	"github.com/minio/minio-go/v7"
)

type MediaRepository interface {
	UploadVideo(file *entity.MediaFileMeta, userID string) (entity.MediaUploadInfo, error)
	UploadImage(file *entity.MediaFileMeta, userID string) (entity.MediaUploadInfo, error)
	ConvertStatus(id string) (entity.MediaConvertStatus, error)
}

type Media struct {
	redisClient *redis.Client
	minioClient *minio.Client
	bucketName  string
}

func NewMedia(redisClient *redis.Client, minioClient *minio.Client, bucketName string) MediaRepository {
	return &Media{
		redisClient: redisClient,
		minioClient: minioClient,
		bucketName:  bucketName,
	}
}

func (r Media) UploadVideo(file *entity.MediaFileMeta, userID string) (entity.MediaUploadInfo, error) {
	return entity.MediaUploadInfo{}, nil
}
func (r Media) UploadImage(file *entity.MediaFileMeta, userID string) (entity.MediaUploadInfo, error) {
	return entity.MediaUploadInfo{}, nil
}
func (r Media) ConvertStatus(id string) (entity.MediaConvertStatus, error) {
	return entity.MediaConvertStatus{}, nil
}
