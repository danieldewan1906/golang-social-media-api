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

type commentService struct {
	db                *sql.DB
	commentRepository domains.CommentRepository
	postService       domains.PostService
}

// Create implements domains.CommentService.
func (service *commentService) Create(ctx context.Context, req dto.CommentRequestDto) error {
	tx, err := service.db.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	_, err = service.postService.FindById(ctx, req.PostId)
	if err != nil {
		return errors.New(err.Error())
	}

	comment := domains.Comment{
		UserId:      req.UserId,
		PostId:      req.PostId,
		TextComment: req.TextComment,
		CreatedAt:   sql.NullTime{Valid: true, Time: time.Now()},
	}

	return service.commentRepository.Create(ctx, tx, comment)
}

// Delete implements domains.CommentService.
func (service *commentService) Delete(ctx context.Context, id int) error {
	tx, err := service.db.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	comment, err := service.commentRepository.FindById(ctx, id)
	util.PanicIfError(err)

	if comment.ID == 0 {
		return errors.New("comment not found")
	}

	return service.commentRepository.Delete(ctx, tx, id)
}

// Update implements domains.CommentService.
func (service *commentService) Update(ctx context.Context, id int, req dto.CommentRequestDto) error {
	tx, err := service.db.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	_, err = service.postService.FindById(ctx, req.PostId)
	if err != nil {
		return errors.New(err.Error())
	}

	comment, err := service.commentRepository.FindById(ctx, id)
	util.PanicIfError(err)

	if comment.ID == 0 {
		return errors.New("comment not found")
	}

	comment.TextComment = req.TextComment
	return service.commentRepository.Update(ctx, tx, comment)
}

func NewCommentService(db *sql.DB, commentRepository domains.CommentRepository, postService domains.PostService) domains.CommentService {
	return &commentService{
		db:                db,
		commentRepository: commentRepository,
		postService:       postService,
	}
}
