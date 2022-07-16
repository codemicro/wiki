package endpoints

import "github.com/gofiber/fiber/v2"

func (e *Endpoints) Index(ctx *fiber.Ctx) error {
	return ctx.SendString("Hello world!")
}
