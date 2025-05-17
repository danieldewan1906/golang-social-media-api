package api

import (
	"context"
	"golang-restful-api/domains"
	"golang-restful-api/dto"
	"golang-restful-api/internal/util"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type authApi struct {
	authService domains.AuthService
}

func NewAuth(router fiber.Router, authService domains.AuthService) {
	authApi := authApi{
		authService: authService,
	}

	auth := router.Group("/auth")

	auth.Post("/login", authApi.Login)
	auth.Post("/register", authApi.Register)
}

func (auth authApi) Login(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}

	res, err := auth.authService.Login(c, req)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(dto.WebResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return ctx.Status(http.StatusOK).JSON(dto.WebResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    res,
	})
}

func (auth authApi) Register(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.RegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}

	fails := util.Validate(req)
	if len(fails) > 0 {
		return ctx.Status(http.StatusBadRequest).JSON(dto.WebResponse{
			Status:  http.StatusBadRequest,
			Message: "Bad Request",
			Data:    fails,
		})
	}

	err := auth.authService.Register(c, req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.WebResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return ctx.Status(http.StatusCreated).JSON(dto.WebResponse{
		Status:  http.StatusCreated,
		Message: "Success",
		Data:    "",
	})
}
