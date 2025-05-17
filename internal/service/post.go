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

type postService struct {
	DB                *sql.DB
	postRepository    domains.PostRepository
	likeRepository    domains.LikeRepository
	commentRepository domains.CommentRepository
}

// FindByUserId implements domains.PostService.
func (service *postService) FindByUserId(ctx context.Context, userId int) ([]dto.PostDto, error) {
	posts, err := service.postRepository.FindByUserId(ctx, userId)
	util.PanicIfError(err)
	var result []dto.PostDto

	if len(posts) < 1 {
		return result, nil
	}

	for _, post := range posts {
		likes, err := service.likeRepository.FindByPostId(ctx, int(post.ID))
		util.PanicIfError(err)

		comments, err := service.commentRepository.FindByPostId(ctx, int(post.ID))
		util.PanicIfError(err)

		result = append(result, dto.PostDto{
			ID:        post.ID,
			UserId:    post.UserId,
			Content:   post.Content.String,
			ImageUrl:  post.ImageUrl.String,
			IsActive:  post.IsActive,
			CreatedAt: post.CreatedAt.Time.String(),
			UpdatedAt: post.UpdatedAt.Time.String(),
			Likes:     likeToLikeDto(likes),
			Comments:  commentToCommentDto(comments),
		})
	}

	return result, nil
}

// ArchievePost implements domains.PostService.
func (service *postService) ArchievePost(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	post, err := service.postRepository.FindById(ctx, id)
	util.PanicIfError(err)

	if post.ID == 0 {
		return errors.New("post not found")
	}

	return service.postRepository.ArchievePost(ctx, tx, id)
}

// Create implements domains.PostService.
func (service *postService) Create(ctx context.Context, req dto.PostRequestDto) error {
	tx, err := service.DB.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	post := domains.Post{
		UserId:    req.UserId,
		Content:   sql.NullString{Valid: true, String: req.Content},
		ImageUrl:  sql.NullString{Valid: true, String: req.Filename},
		IsActive:  true,
		CreatedAt: sql.NullTime{Valid: true, Time: time.Now()},
	}
	err = service.postRepository.Create(ctx, tx, post)
	return err
}

// Delete implements domains.PostService.
func (service *postService) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	post, err := service.postRepository.FindById(ctx, id)
	util.PanicIfError(err)

	if post.ID == 0 {
		return errors.New("post not found")
	}

	return service.postRepository.Delete(ctx, tx, id)
}

// FindAll implements domains.PostService.
func (service *postService) FindAll(ctx context.Context) ([]dto.PostDto, error) {
	posts, err := service.postRepository.FindAll(ctx)
	util.PanicIfError(err)

	var result []dto.PostDto
	for _, post := range posts {
		result = append(result, dto.PostDto{
			ID:        post.ID,
			UserId:    post.UserId,
			Content:   post.Content.String,
			ImageUrl:  post.ImageUrl.String,
			IsActive:  post.IsActive,
			CreatedAt: post.CreatedAt.Time.String(),
			UpdatedAt: post.UpdatedAt.Time.String(),
		})
	}

	return result, nil
}

// FindById implements domains.PostService.
func (service *postService) FindById(ctx context.Context, id int) (dto.PostDto, error) {
	post, err := service.postRepository.FindById(ctx, id)
	util.PanicIfError(err)

	if post.ID == 0 {
		return dto.PostDto{}, errors.New("post not found")
	}

	likes, err := service.likeRepository.FindByPostId(ctx, id)
	util.PanicIfError(err)

	comments, err := service.commentRepository.FindByPostId(ctx, id)
	util.PanicIfError(err)

	return dto.PostDto{
		ID:        post.ID,
		UserId:    post.UserId,
		Content:   post.Content.String,
		ImageUrl:  post.ImageUrl.String,
		IsActive:  post.IsActive,
		CreatedAt: post.CreatedAt.Time.String(),
		UpdatedAt: post.UpdatedAt.Time.String(),
		Likes:     likeToLikeDto(likes),
		Comments:  commentToCommentDto(comments),
	}, nil
}

// Update implements domains.PostService.
func (service *postService) Update(ctx context.Context, id int, req dto.PostRequestDto) error {
	tx, err := service.DB.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	post, _ := service.postRepository.FindByUserIdAndId(ctx, req.UserId, id)
	if post.ID == 0 {
		return errors.New("post not found")
	}

	reqPost := domains.Post{
		ID:        post.ID,
		UserId:    req.UserId,
		Content:   sql.NullString{Valid: true, String: req.Content},
		ImageUrl:  sql.NullString{Valid: true, String: req.Filename},
		IsActive:  true,
		UpdatedAt: sql.NullTime{Valid: true, Time: time.Now()},
	}
	err = service.postRepository.Update(ctx, tx, reqPost)
	return err
}

func NewPostService(db *sql.DB, postRepository domains.PostRepository, likeRepository domains.LikeRepository, commentRepository domains.CommentRepository) domains.PostService {
	return &postService{
		DB:                db,
		postRepository:    postRepository,
		likeRepository:    likeRepository,
		commentRepository: commentRepository,
	}
}

func likeToLikeDto(like []domains.Like) []dto.LikeDto {
	var likeDto []dto.LikeDto
	for _, v := range like {
		likeDto = append(likeDto, dto.LikeDto{
			UserId:    v.UserId,
			CreatedAt: v.CreatedAt.Time.String(),
		})
	}
	return likeDto
}

func commentToCommentDto(comment []domains.Comment) []dto.CommentDto {
	var commentDto []dto.CommentDto
	for _, v := range comment {
		commentDto = append(commentDto, dto.CommentDto{
			ID:          int(v.ID),
			UserId:      v.UserId,
			PostId:      v.PostId,
			TextComment: v.TextComment,
			CreatedAt:   v.CreatedAt.Time.String(),
		})
	}
	return commentDto
}
