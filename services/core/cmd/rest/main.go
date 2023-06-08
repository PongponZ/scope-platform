package main

import (
	"github.com/PongponZ/scope-platform/core/config"
	"github.com/gofiber/fiber/v2"
)

func main() {
	conf := config.GetConfig()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello world")
	})

	app.Listen(conf.ServerPort)
}
