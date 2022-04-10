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

	if req.Title == "" || req.Body == "" || req.Category < 0 || req.Category > 10 {
		return c.Status(http.StatusBadRequest).SendString("Http Body is not Matched Fakenews Request Scheme")
	}

	uuid, err := requestSvc.CreateAnalyze(c.Context(), req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(http.StatusOK).SendString(uuid)
}

func GetAnalyze(c *fiber.Ctx) error {
	requestId := c.Params("id", "")
	if requestId == "" {
		return c.Status(http.StatusBadRequest).SendString("Id Should no be empty")
	}

	request, err := requestSvc.GetRequest(c.Context(), requestId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(http.StatusOK).JSON(request)
}

func ListAnalyze(c *fiber.Ctx) error {
	requests, err := requestSvc.ListRequests(c.Context())
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(http.StatusOK).JSON(requests)
}
