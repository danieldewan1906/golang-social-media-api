package domains

import (
	"context"
	"golang-restful-api/dto"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) error
	Login(ctx context.Context, req dto.LoginRequest) (dto.AuthResponse, error)
}
