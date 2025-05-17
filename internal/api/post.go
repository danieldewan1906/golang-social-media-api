package api

import (
	"context"
	"fmt"
	"golang-restful-api/domains"
	"golang-restful-api/dto"
	"golang-restful-api/internal/config"
	"golang-restful-api/internal/middleware"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type postApi struct {
	postService domains.PostService
	config      config.Config
}

func NewPostApi(router fiber.Router, config *config.Config, postService domains.PostService) {
	postApi := postApi{
		postService: postService,
		config:      *config,
	}

	post := router.Group("/post")

	post.Get("all", middleware.JwtWithRolePermission("USER"), postApi.FindAll)
	post.Get(":id", middleware.JwtWithRolePermission("USER"), postApi.FindById)
	post.Post("", middleware.JwtWithRolePermission("USER"), postApi.Create)
	post.Put(":id", middleware.JwtValidateUser(), middleware.JwtWithRolePermission("USER"), postApi.Update)
	post.Delete("", middleware.JwtValidateUser(), middleware.JwtWithRolePermission("USER"), postApi.DeleteArchivePost)
}

func (post postApi) FindAll(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	result, err := post.postService.FindAll(c)
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
		Data:    result,
	})
}

func (post postApi) FindById(ctx *fiber.Ctx) error {
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

	result, err := post.postService.FindById(c, id)
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
		Data:    result,
	})
}

func (post postApi) Create(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	file, err := ctx.FormFile("file")
	var extension string
	var filename string
	if err == nil {
		extension = filepath.Ext(file.Filename)
		filename = strings.TrimSuffix(file.Filename, extension) + "_" + uuid.NewString()
		savePath := fmt.Sprintf("%s/%s", post.config.FileUpload.Path, filename+extension)
		if err := ctx.SaveFile(file, savePath); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(dto.WebResponse{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
				Data:    nil,
			})
		}
	}

	content := ctx.FormValue("content")
	userId := int(ctx.Locals("userId").(float64))
	postRequest := dto.PostRequestDto{
		UserId:   userId,
		Content:  content,
		Filename: "/" + filename + extension,
	}

	err = post.postService.Create(c, postRequest)
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

func (post postApi) Update(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	file, err := ctx.FormFile("file")
	var extension string
	var filename string
	if err == nil {
		extension = filepath.Ext(file.Filename)
		filename = strings.TrimSuffix(file.Filename, extension) + "_" + uuid.NewString()
		savePath := fmt.Sprintf("%s/%s", post.config.FileUpload.Path, filename+extension)
		if err := ctx.SaveFile(file, savePath); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(dto.WebResponse{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
				Data:    nil,
			})
		}
	}

	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.WebResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	content := ctx.FormValue("content")
	userId := int(ctx.Locals("userId").(float64))
	postRequest := dto.PostRequestDto{
		UserId:   userId,
		Content:  content,
		Filename: "/" + filename + extension,
	}

	err = post.postService.Update(c, id, postRequest)
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

func (post postApi) DeleteArchivePost(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var request dto.PostDeleteRequestDto
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(dto.WebResponse{
			Status:  http.StatusUnprocessableEntity,
			Message: err.Error(),
			Data:    nil,
		})
	}

	var err error
	if request.IsDeletePost {
		err = post.postService.Delete(c, request.ID)
	} else {
		err = post.postService.ArchievePost(c, request.ID)
	}

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
