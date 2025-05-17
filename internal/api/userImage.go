package api

import (
	"context"
	"fmt"
	"golang-restful-api/domains"
	"golang-restful-api/dto"
	"golang-restful-api/internal/config"
	"golang-restful-api/internal/middleware"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type userImageApi struct {
	userImageService domains.UserImagesService
	config           config.Config
}

func NewUserImages(router fiber.Router, cnf *config.Config, userImageService domains.UserImagesService) {
	userImageApi := userImageApi{
		userImageService: userImageService,
		config:           *cnf,
	}

	userImages := router.Group("/user-image/:userId")

	userImages.Get("", middleware.JwtWithRolePermission("USER", "ADMIN"), userImageApi.FindByUserId)
	userImages.Post("", middleware.JwtWithRolePermission("USER"), userImageApi.Create)
	userImages.Put("", middleware.JwtWithRolePermission("USER"), userImageApi.Update)
	userImages.Delete("", middleware.JwtWithRolePermission("USER"), userImageApi.Delete)
}

func (api userImageApi) FindByUserId(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	userId, err := strconv.Atoi(ctx.Params("userId"))
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.WebResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	result, err := api.userImageService.FindByUserId(c, userId)
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

func (api userImageApi) Create(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.WebResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	userId, err := strconv.Atoi(ctx.Params("userId"))
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.WebResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	usrImg, _ := api.userImageService.FindByUserId(c, userId)
	if usrImg.ID != 0 {
		return ctx.Status(http.StatusBadRequest).JSON(dto.WebResponse{
			Status:  http.StatusBadRequest,
			Message: "Image already set!",
			Data:    nil,
		})
	}

	extension := filepath.Ext(file.Filename)
	filename := strings.TrimSuffix(file.Filename, extension) + "_" + uuid.NewString()
	savePath := fmt.Sprintf("%s/%s", api.config.FileUpload.Path, filename+extension)
	if err := ctx.SaveFile(file, savePath); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.WebResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
	}

	payload := dto.UserImageRequestDto{
		Filename:  filename,
		Extension: extension,
	}
	err = api.userImageService.Save(c, payload, userId)
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

func (api userImageApi) Update(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.WebResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	userId, err := strconv.Atoi(ctx.Params("userId"))
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.WebResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	usrImg, err := api.userImageService.FindByUserId(c, userId)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.WebResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	if err := os.Remove(api.config.FileUpload.Path + usrImg.ImageUrl + usrImg.Extension); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.WebResponse{
			Status:  http.StatusInternalServerError,
			Message: "Gagal hapus file",
			Data:    err.Error(),
		})
	}

	extension := filepath.Ext(file.Filename)
	filename := strings.TrimSuffix(file.Filename, extension) + "_" + uuid.NewString()
	savePath := fmt.Sprintf("%s/%s", api.config.FileUpload.Path, filename+extension)
	if err := ctx.SaveFile(file, savePath); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.WebResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
	}

	payload := dto.UserImageRequestDto{
		Filename:  filename,
		Extension: extension,
	}
	err = api.userImageService.Update(c, payload, userId)
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

func (api userImageApi) Delete(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	userId, err := strconv.Atoi(ctx.Params("userId"))
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.WebResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	usrImg, _ := api.userImageService.FindByUserId(c, userId)
	if err := os.Remove(api.config.FileUpload.Path + usrImg.ImageUrl + usrImg.Extension); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.WebResponse{
			Status:  http.StatusInternalServerError,
			Message: "Gagal hapus file",
			Data:    nil,
		})
	}

	err = api.userImageService.Delete(c, userId)
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
