package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/PongponZ/scope-platform/core/internal/entity"
	"github.com/PongponZ/scope-platform/core/internal/repository"
	"github.com/PongponZ/scope-platform/core/pkg/common/array"
	"github.com/PongponZ/scope-platform/core/pkg/rabbitmq"
)

type MediaUsecase interface {
	IsAllowType(meta *entity.MediaFileMeta) error
	GetConvertStatus(id string) (entity.MediaConvertStatus, error)
	GetMetaFromBuffer(file *multipart.FileHeader, buffer *multipart.File) *entity.MediaFileMeta
	UploadVideo(meta *entity.MediaFileMeta, userID string) (*entity.MediaUploaded, error)
}

type Media struct {
	mediaDomain       string
	mediaRepository   repository.MediaRepository
	messageBroker     *rabbitmq.RabbitMQ
	mediaConvertQueue string
	allowFileType     []string
}

func NewMedia(mediaDomain string, mediaRepo repository.MediaRepository, broker *rabbitmq.RabbitMQ, convertQueue string, allowFileType []string) MediaUsecase {
	return &Media{
		mediaDomain:       mediaDomain,
		mediaRepository:   mediaRepo,
		messageBroker:     broker,
		mediaConvertQueue: convertQueue,
		allowFileType:     allowFileType,
	}
}

func (u Media) IsAllowType(meta *entity.MediaFileMeta) error {
	if found := array.Contains(u.allowFileType, meta.Extension); !found {
		return errors.New(ErrorInvalidMediaType)
	}
	return nil
}

func (u Media) GetMetaFromBuffer(file *multipart.FileHeader, buffer *multipart.File) *entity.MediaFileMeta {
	typeSlice := strings.Split(file.Header["Content-Type"][0], "/")
	return &entity.MediaFileMeta{
		FileName:    file.Filename,
		Buffer:      buffer,
		ContentType: file.Header["Content-Type"][0],
		Extension:   typeSlice[1],
		Type:        typeSlice[0],
		Size:        file.Size,
	}
}

func (u Media) GetConvertStatus(id string) (entity.MediaConvertStatus, error) {
	return u.mediaRepository.ConvertStatus(id)
}

func (u Media) UploadVideo(meta *entity.MediaFileMeta, userID string) (*entity.MediaUploaded, error) {
	info, err := u.mediaRepository.UploadVideo(meta, userID)
	if err != nil {
		return nil, err
	}

	payload := u.UploadInfoToMediaConvertJob(info)
	jsonPayload, _ := json.Marshal(payload)

	if err = u.messageBroker.Publish(u.mediaConvertQueue, jsonPayload); err != nil {
		return nil, errors.New(ErrorCannotPublishMediaConvertMessage)
	}

	endpoint := u.CreateMediaEndPoint(info.ID+".m3u8", userID)
	return &entity.MediaUploaded{
		ID:        info.ID,
		Endpoint:  endpoint,
		Timestamp: time.Now().Unix(),
	}, nil
}

func (u Media) CreateMediaEndPoint(filname string, userID string) string {
	return fmt.Sprintf("%s/%s/%s", u.mediaDomain, userID, filname)
}

func (u Media) UploadInfoToMediaConvertJob(info entity.MediaUploadInfo) entity.MediaConvertPayload {
	return entity.MediaConvertPayload{
		ID:         info.ID,
		FileName:   info.Filename,
		Path:       fmt.Sprintf("%s/%s/%s", info.UserId, info.Type, info.Filename),
		Output:     info.ID,
		OutputPath: fmt.Sprintf("%s/%s", info.UserId, info.Type),
		Type:       info.Type,
	}
}
