package service

import (
	"context"
	"database/sql"
	"errors"
	"golang-restful-api/domains"
	"golang-restful-api/dto"
	"golang-restful-api/internal/config"
	"golang-restful-api/internal/util"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	DB                   *sql.DB
	conf                 *config.Config
	userRepository       domains.UserRepository
	userDetailRepository domains.UserDetailRepository
}

func NewAuth(DB *sql.DB, cnf *config.Config, userRepository domains.UserRepository, userDetailRepository domains.UserDetailRepository) domains.AuthService {
	return authService{
		DB:                   DB,
		conf:                 cnf,
		userRepository:       userRepository,
		userDetailRepository: userDetailRepository,
	}
}

// Login implements domains.AuthService.
func (a authService) Login(ctx context.Context, req dto.LoginRequest) (dto.AuthResponse, error) {
	tx, err := a.DB.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	user, err := a.userRepository.FindByEmailUsername(ctx, req.Email, req.Username)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	if user.ID == 0 {
		return dto.AuthResponse{}, errors.New("user not found or deactivate account! please register again")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return dto.AuthResponse{}, errors.New("unauthorized")
	}

	claim := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Duration(a.conf.Jwt.Exp) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claim)
	tokenStr, err := token.SignedString([]byte(a.conf.Jwt.Key))
	if err != nil {
		return dto.AuthResponse{}, errors.New("unauthorized")
	}

	user.Token = sql.NullString{Valid: true, String: tokenStr}
	err = a.userRepository.Update(ctx, tx, &user)
	util.PanicIfError(err)

	return dto.AuthResponse{
		Token: tokenStr,
	}, nil
}

// Register implements domains.AuthService.
func (a authService) Register(ctx context.Context, req dto.RegisterRequest) error {
	tx, err := a.DB.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	user, err := a.userRepository.FindByEmailUsername(ctx, req.Email, req.Username)
	util.PanicIfError(err)

	if user.ID != 0 {
		return errors.New("user already registered")
	}

	if req.Password != req.Repassword {
		return errors.New("wrong password")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		util.PanicIfError(errors.New("error bcrypt"))
	}

	reqUser := domains.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(password),
		IsActive:  true,
		CreatedAt: sql.NullTime{Valid: true, Time: time.Now()},
		Role:      req.Role,
	}

	result, err := a.userRepository.Save(ctx, tx, &reqUser)
	util.PanicIfError(err)

	var reqUserDetail = domains.UserDetail{
		UserId:    int(result.ID),
		FirstName: req.FirstName,
		LastName:  sql.NullString{Valid: true, String: req.LastName},
		CreatedAt: sql.NullTime{Valid: true, Time: time.Now()},
	}

	err = a.userDetailRepository.Save(ctx, tx, &reqUserDetail)
	util.PanicIfError(err)
	return err
}
