package api

import (
	"context"
	"golang-restful-api/domains"
	"golang-restful-api/dto"
	"golang-restful-api/internal/middleware"
	"golang-restful-api/internal/util"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type userDetailApi struct {
	userDetailService domains.UserDetailService
}

func NewUserDetail(router fiber.Router, userDetailService domains.UserDetailService) {
	userDetailApi := userDetailApi{
		userDetailService: userDetailService,
	}

	userDetail := router.Group("/user")

	userDetail.Get("current", middleware.JwtWithRolePermission("ADMIN", "USER"), userDetailApi.CurrentUser)
	userDetail.Get("", middleware.JwtWithRolePermission("ADMIN"), userDetailApi.FindAll)
	userDetail.Put("", middleware.JwtWithRolePermission("USER"), userDetailApi.Update)
	userDetail.Delete(":userId", middleware.JwtWithRolePermission("ADMIN", "USER"), userDetailApi.InActiveUser)
}

func (userDetail userDetailApi) CurrentUser(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	userId := int(ctx.Locals("userId").(float64))
	customer, err := userDetail.userDetailService.FindByUserId(c, userId)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.WebResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	if customer.ID == 0 {
		return ctx.Status(http.StatusNotFound).
			JSON(dto.WebResponse{
				Status:  http.StatusNotFound,
				Message: "User not found",
				Data:    nil,
			})
	}

	return ctx.JSON(dto.WebResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    customer,
	})
}

func (userDetail userDetailApi) FindAll(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.UserDetailRequestDto
	if err := ctx.QueryParser(&req); err != nil {
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}

	res, err := userDetail.userDetailService.FindAll(c, req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.WebResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return ctx.JSON(dto.WebResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    res,
	})
}

func (userDetail userDetailApi) Update(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var request dto.UpdateUserRequestDto
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

	err := userDetail.userDetailService.Update(c, request)
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

func (userDetail userDetailApi) InActiveUser(ctx *fiber.Ctx) error {
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

	err = userDetail.userDetailService.InActiveUser(c, userId)
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
