package api

import (
	"context"
	"golang-restful-api/domains"
	"golang-restful-api/dto"
	"golang-restful-api/internal/middleware"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type likeApi struct {
	likeService domains.LikeService
}

func NewLikeApi(router fiber.Router, likeService domains.LikeService) {
	likeApi := likeApi{
		likeService: likeService,
	}

	like := router.Group("/like")

	like.Post("", middleware.JwtValidateUser(), middleware.JwtWithRolePermission("USER"), likeApi.Create)
	like.Delete("", middleware.JwtValidateUser(), middleware.JwtWithRolePermission("USER"), likeApi.Delete)
}

func (like likeApi) Create(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var request dto.LikeRequestDto
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(dto.WebResponse{
			Status:  http.StatusUnprocessableEntity,
			Message: err.Error(),
			Data:    nil,
		})
	}

	err := like.likeService.Create(c, request)
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
		Data:    "",
	})
}

func (like likeApi) Delete(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var request dto.LikeRequestDto
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(dto.WebResponse{
			Status:  http.StatusUnprocessableEntity,
			Message: err.Error(),
			Data:    nil,
		})
	}

	err := like.likeService.Delete(c, request)
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
		Data:    "",
	})
}
