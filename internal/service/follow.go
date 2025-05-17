package service

import (
	"context"
	"database/sql"
	"errors"
	"golang-restful-api/domains"
	"golang-restful-api/dto"
	"golang-restful-api/internal/util"
	"time"
)

type followService struct {
	db               *sql.DB
	followRepository domains.FollowRepository
	userService      domains.UserDetailService
}

// Create implements domains.FollowService.
func (service *followService) Create(ctx context.Context, req dto.FollowRequestDto) error {
	tx, err := service.db.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	user, err := service.userService.FindByUserId(ctx, req.FollowingId)
	if err != nil {
		return errors.New(err.Error())
	}

	if user.ID == 0 {
		return errors.New("following user not found")
	}

	follow, err := service.followRepository.FindByFollowingAndFollowerId(ctx, req.UserId, req.FollowingId)
	util.PanicIfError(err)

	if follow.FollowingId != 0 {
		return errors.New("user has been followed by you")
	}

	follow.FollowerId = req.UserId
	follow.FollowingId = req.FollowingId
	follow.CreatedAt = sql.NullTime{Valid: true, Time: time.Now()}
	return service.followRepository.Create(ctx, tx, follow)
}

// Delete implements domains.FollowService.
func (service *followService) Delete(ctx context.Context, req dto.FollowRequestDto) error {
	tx, err := service.db.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	user, err := service.userService.FindByUserId(ctx, req.FollowingId)
	if err != nil {
		return errors.New(err.Error())
	}

	if user.ID == 0 {
		return errors.New("following user not found")
	}

	follow, err := service.followRepository.FindByFollowingAndFollowerId(ctx, req.UserId, req.FollowingId)
	util.PanicIfError(err)

	if follow.FollowingId == 0 && follow.FollowerId == 0 {
		return errors.New("users you haven't followed yet")
	}

	return service.followRepository.Delete(ctx, tx, follow)
}

func NewFollowService(db *sql.DB, followRepository domains.FollowRepository, userService domains.UserDetailService) domains.FollowService {
	return &followService{
		db:               db,
		followRepository: followRepository,
		userService:      userService,
	}
}
