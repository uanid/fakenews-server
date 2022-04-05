package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/uanid/fakenews-server/pkg/services"
	"github.com/uanid/fakenews-server/pkg/types"
)

var requestSvc *services.RequestService

func SetRequestSvc(svc *services.RequestService) {
	requestSvc = svc
}

func Ping(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).SendString("Pong")
}

func RequestAnalyze(c *fiber.Ctx) error {
	req := &types.FakeNewsReq{}
	err := c.BodyParser(req)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Invalid Request Struct")
	}

	uuid, err := requestSvc.CreateRequest(c.Context(), req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(http.StatusOK).SendString("uuid:" + uuid)
}

func GetAnalyze(c *fiber.Ctx) error {
	return nil
}
