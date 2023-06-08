package handlers

import "github.com/gofiber/fiber/v2"

var (
	ErrCodeBadRequest          string = "bad_request"
	ErrCodeNotFound            string = "not_found"
	ErrCodePublishMessageQueue string = "publish_message_queue"
	ErrCodeInternalError       string = "internal_error"
)

func NotFound(c *fiber.Ctx, message string) error {
	errResponse := NewErrorResponse(ErrCodeNotFound, message)
	return c.Status(fiber.StatusNotFound).JSON(errResponse)
}

func BadRequest(c *fiber.Ctx, message string) error {
	errResponse := NewErrorResponse(ErrCodeBadRequest, message)
	return c.Status(fiber.StatusBadRequest).JSON(errResponse)
}

func InternalError(c *fiber.Ctx, code string, message string) error {
	errResponse := NewErrorResponse(code, message)
	return c.Status(fiber.StatusInternalServerError).JSON(errResponse)
}

func NewErrorResponse(code string, message string) Response {
	return Response{
		Success:   false,
		ErrorCode: code,
		Message:   message,
	}
}
