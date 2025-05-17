package api

import (
	"context"
	"golang-restful-api/domains"
	"golang-restful-api/dto"
	"golang-restful-api/internal/middleware"
	"golang-restful-api/internal/util"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type followApi struct {
	followService domains.FollowService
}

func NewFollowApi(router fiber.Router, followService domains.FollowService) {
	followApi := followApi{
		followService: followService,
	}

	follow := router.Group("/follow")
	follow.Post("", middleware.JwtValidateUser(), middleware.JwtWithRolePermission("USER"), followApi.Create)
	follow.Delete("", middleware.JwtValidateUser(), middleware.JwtWithRolePermission("USER"), followApi.Delete)
}

func (follow followApi) Create(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var request dto.FollowRequestDto
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(dto.WebResponse{
			Status:  http.StatusUnprocessableEntity,
			Message: err.Error(),
			Data:    nil,
		})
	}

	validate := util.Validate(request)
	if len(validate) > 0 {
		return ctx.Status(http.StatusBadRequest).JSON(dto.WebResponse{
			Status:  http.StatusBadRequest,
			Message: "Bad Request",
			Data:    validate,
		})
	}

	err := follow.followService.Create(c, request)
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
		Data:    nil,
	})
}

func (follow followApi) Delete(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var request dto.FollowRequestDto
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(dto.WebResponse{
			Status:  http.StatusUnprocessableEntity,
			Message: err.Error(),
			Data:    nil,
		})
	}

	validate := util.Validate(request)
	if len(validate) > 0 {
		return ctx.Status(http.StatusBadRequest).JSON(dto.WebResponse{
			Status:  http.StatusBadRequest,
			Message: "Bad Request",
			Data:    validate,
		})
	}

	err := follow.followService.Delete(c, request)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.WebResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return ctx.Status(http.StatusOK).JSON(dto.WebResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    nil,
	})
}
