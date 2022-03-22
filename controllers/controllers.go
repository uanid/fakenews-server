package controllers

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func Ping(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).SendString("Pong")
}

func RequestAnalyze(c *fiber.Ctx) error  {
	return nil
}

func GetAnalyze(c *fiber.Ctx) error  {
	return nil
}
