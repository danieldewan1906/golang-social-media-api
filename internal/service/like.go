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

type likeService struct {
	db             *sql.DB
	likeRepository domains.LikeRepository
	postService    domains.PostService
}

// Create implements domains.LikeService.
func (service *likeService) Create(ctx context.Context, req dto.LikeRequestDto) error {
	tx, err := service.db.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	data, err := service.likeRepository.FindByUserIdAndPostId(ctx, req.UserId, req.PostId)
	if err != nil {
		return errors.New(err.Error())
	}

	if data.ID != 0 {
		return errors.New("post ini sudah disukai")
	}

	_, err = service.postService.FindById(ctx, req.PostId)
	if err != nil {
		return errors.New(err.Error())
	}

	like := domains.Like{
		UserId:    req.UserId,
		PostId:    req.PostId,
		CreatedAt: sql.NullTime{Valid: true, Time: time.Now()},
	}

	return service.likeRepository.Create(ctx, tx, like)
}

// Delete implements domains.LikeService.
func (service *likeService) Delete(ctx context.Context, req dto.LikeRequestDto) error {
	tx, err := service.db.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	_, err = service.postService.FindById(ctx, req.PostId)
	if err != nil {
		return errors.New(err.Error())
	}

	return service.likeRepository.Delete(ctx, tx, req.UserId, req.PostId)
}

func NewLikeService(db *sql.DB, likeRepository domains.LikeRepository, postService domains.PostService) domains.LikeService {
	return &likeService{
		db:             db,
		likeRepository: likeRepository,
		postService:    postService,
	}
}
