package api

import (
	"context"
	"golang-restful-api/domains"
	"golang-restful-api/dto"
	"golang-restful-api/internal/middleware"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type commentApi struct {
	commentService domains.CommentService
}

func NewCommentApi(router fiber.Router, commentService domains.CommentService) {
	commentApi := commentApi{
		commentService: commentService,
	}

	comment := router.Group("/comment")

	comment.Post("", middleware.JwtValidateUser(), middleware.JwtWithRolePermission("USER"), commentApi.Create)
	comment.Put(":id", middleware.JwtValidateUser(), middleware.JwtWithRolePermission("USER"), commentApi.Update)
	comment.Delete(":id", middleware.JwtWithRolePermission("USER"), commentApi.Delete)
}

func (comment commentApi) Create(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var request dto.CommentRequestDto
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(dto.WebResponse{
			Status:  http.StatusUnprocessableEntity,
			Message: err.Error(),
			Data:    nil,
		})
	}

	err := comment.commentService.Create(c, request)
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

func (comment commentApi) Update(ctx *fiber.Ctx) error {
	log.Println("MASUK UPDATE COMMENT")
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var request dto.CommentRequestDto
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(dto.WebResponse{
			Status:  http.StatusUnprocessableEntity,
			Message: err.Error(),
			Data:    nil,
		})
	}

	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.WebResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	err = comment.commentService.Update(c, id, request)
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

func (comment commentApi) Delete(ctx *fiber.Ctx) error {
	log.Println("MASUK DELETE COMMENT")
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.WebResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	err = comment.commentService.Delete(c, id)
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
