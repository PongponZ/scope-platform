package handlers

import (
	"github.com/PongponZ/scope-platform/core/internal/errors"
	"github.com/PongponZ/scope-platform/core/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type MediaHandler interface {
	UploadVideo(c *fiber.Ctx) error
	ConvertStatus(c *fiber.Ctx) error
}

type Media struct {
	mediaUsecase usecase.MediaUsecase
}

func NewMedia(mediaUsecase usecase.MediaUsecase) MediaHandler {
	return Media{
		mediaUsecase: mediaUsecase,
	}
}

func (h Media) ConvertStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return BadRequest(c, "missing id")
	}

	result, err := h.mediaUsecase.GetConvertStatus(id)
	if err != nil {
		return InternalError(c, ErrCodeInternalError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: "get media converting success",
		Data:    result,
	})
}

func (h Media) UploadVideo(c *fiber.Ctx) error {
	file, err := c.FormFile("media")
	if err != nil {
		return BadRequest(c, err.Error())
	}

	buffer, err := file.Open()
	if err != nil {
		return BadRequest(c, err.Error())
	}
	defer buffer.Close()

	metaFile := h.mediaUsecase.GetMetaFromBuffer(file, &buffer)
	if err := h.mediaUsecase.IsAllowType(metaFile); err != nil {
		return BadRequest(c, err.Error())
	}

	result, err := h.mediaUsecase.UploadVideo(metaFile, c.Locals("userID").(string))
	if err != nil {
		if err.Error() == errors.ErrorCannotPublishMediaConvertMessage {
			return InternalError(c, ErrCodePublishMessageQueue, err.Error())
		}
		return InternalError(c, ErrCodeInternalError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Message: "upload video success",
		Data:    result,
	})
}
