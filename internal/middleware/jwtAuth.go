package middleware

import (
	"golang-restful-api/dto"
	"golang-restful-api/internal/config"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserIdStruct struct {
	UserId int `json:"userId" param:"userId" form:"userId" query:"userId"`
}

func JwtWithRolePermission(roles ...string) fiber.Handler {
	config := config.Get()
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(http.StatusUnauthorized).JSON(dto.WebResponse{
				Status:  http.StatusUnauthorized,
				Message: "UNAUTHORIZED",
				Data:    nil,
			})
		}

		tokenStr := strings.Split(authHeader, "Bearer ")[1]
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.Jwt.Key), nil
		})

		if err != nil || !token.Valid {
			return c.Status(http.StatusUnauthorized).JSON(dto.WebResponse{
				Status:  http.StatusUnauthorized,
				Message: "UNAUTHORIZED",
				Data:    nil,
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(http.StatusUnauthorized).JSON(dto.WebResponse{
				Status:  http.StatusUnauthorized,
				Message: "UNAUTHORIZED",
				Data:    nil,
			})
		}

		role := claims["role"]
		email := claims["email"]
		c.Locals("userRole", role)
		c.Locals("email", email)
		c.Locals("userId", claims["id"])

		for _, r := range roles {
			if role == r {
				return c.Next()
			}
		}

		return c.Status(http.StatusForbidden).JSON(dto.WebResponse{
			Status:  http.StatusForbidden,
			Message: "FORBIDDEN",
			Data:    nil,
		})
	}
}

func JwtValidateUser() fiber.Handler {
	config := config.Get()
	return func(c *fiber.Ctx) error {
		var body UserIdStruct
		if err := c.BodyParser(&body); err != nil {
			body.UserId = 0
		}

		// user
		userIdBody := strconv.Itoa(body.UserId)
		userIdQuery := c.Query("userId")
		userIdForm := c.FormValue("userId")
		userIdParam := c.Params("userId")

		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(http.StatusUnauthorized).JSON(dto.WebResponse{
				Status:  http.StatusUnauthorized,
				Message: "UNAUTHORIZED",
				Data:    nil,
			})
		}

		tokenStr := strings.Split(authHeader, "Bearer ")[1]
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.Jwt.Key), nil
		})

		if err != nil || !token.Valid {
			return c.Status(http.StatusUnauthorized).JSON(dto.WebResponse{
				Status:  http.StatusUnauthorized,
				Message: "UNAUTHORIZED",
				Data:    nil,
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(http.StatusUnauthorized).JSON(dto.WebResponse{
				Status:  http.StatusUnauthorized,
				Message: "UNAUTHORIZED",
				Data:    nil,
			})
		}

		id := strconv.Itoa(int(claims["id"].(float64)))
		role := claims["role"]
		email := claims["email"]
		c.Locals("userRole", role)
		c.Locals("email", email)
		c.Locals("userId", id)

		if userIdBody != "" || userIdForm != "" || userIdParam != "" || userIdQuery != "" {
			if id == userIdBody || id == userIdForm || id == userIdParam || id == userIdQuery {
				return c.Next()
			}
		}

		return c.Status(http.StatusForbidden).JSON(dto.WebResponse{
			Status:  http.StatusForbidden,
			Message: "FORBIDDEN",
			Data:    nil,
		})
	}
}
