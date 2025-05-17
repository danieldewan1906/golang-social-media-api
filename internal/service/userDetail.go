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

type userDetailService struct {
	DB                   *sql.DB
	userDetailRepository domains.UserDetailRepository
	userImageService     domains.UserImagesService
	followRepository     domains.FollowRepository
	postService          domains.PostService
}

func NewUserDetail(DB *sql.DB, userDetailRepository domains.UserDetailRepository, userImageService domains.UserImagesService, followRepository domains.FollowRepository, postService domains.PostService) domains.UserDetailService {
	return &userDetailService{
		DB:                   DB,
		userDetailRepository: userDetailRepository,
		userImageService:     userImageService,
		followRepository:     followRepository,
		postService:          postService,
	}
}

// Index implements domain.CustomerService.
func (service *userDetailService) FindAll(ctx context.Context, req dto.UserDetailRequestDto) ([]dto.UserDetailDto, error) {
	users, err := service.userDetailRepository.FindAll(ctx, req)
	if err != nil {
		return nil, err
	}

	var userDetailDto []dto.UserDetailDto
	for _, v := range users {
		userDetailDto = append(userDetailDto, dto.UserDetailDto{
			ID:        v.UserId,
			FirstName: v.FirstName,
			LastName:  v.LastName.String,
			Address:   v.Address.String,
			CreatedAt: v.CreatedAt.Time.String(),
		})
	}

	return userDetailDto, nil
}

// FindById implements domains.CustomerService.
func (service *userDetailService) FindByUserId(ctx context.Context, userId int) (dto.UserDetailDto, error) {
	user, err := service.userDetailRepository.FindByUserId(ctx, userId)
	if err != nil {
		return dto.UserDetailDto{}, err
	}

	userImage, _ := service.userImageService.FindByUserId(ctx, userId)

	posts, err := service.postService.FindByUserId(ctx, userId)
	util.PanicIfError(err)
	follower, err := service.followRepository.FindByFollowerId(ctx, userId)
	util.PanicIfError(err)
	following, err := service.followRepository.FindByFollowingId(ctx, userId)
	util.PanicIfError(err)

	return dto.UserDetailDto{
		ID:        user.UserId,
		FirstName: user.FirstName,
		LastName:  user.LastName.String,
		Address:   user.Address.String,
		CreatedAt: user.CreatedAt.Time.String(),
		UserImage: &userImage,
		Posts: dto.DetailPostDto{
			Total: len(posts),
			Data:  posts,
		},
		Followers: dto.DetailFollowDto{
			Total: len(following),
			Data:  followToFollowDto(following),
		},
		Following: dto.DetailFollowDto{
			Total: len(follower),
			Data:  followToFollowDto(follower),
		},
	}, nil
}

// Update implements domains.UserDetailService.
func (service *userDetailService) Update(ctx context.Context, req dto.UpdateUserRequestDto) error {
	tx, err := service.DB.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	user, err := service.userDetailRepository.FindByUserId(ctx, req.UserId)
	util.PanicIfError(err)

	if user.ID == 0 {
		return errors.New("user not found")
	}

	requestUser := domains.UserDetail{
		ID:        user.ID,
		UserId:    user.UserId,
		FirstName: req.FirstName,
		LastName:  sql.NullString{Valid: true, String: req.LastName},
		Address:   sql.NullString{Valid: true, String: req.Address},
		UpdatedAt: sql.NullTime{Valid: true, Time: time.Now()},
	}

	return service.userDetailRepository.Update(ctx, tx, &requestUser)
}

// InActiveUser implements domains.UserDetailService.
func (service *userDetailService) InActiveUser(ctx context.Context, userId int) error {
	tx, err := service.DB.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	user, err := service.userDetailRepository.FindByUserId(ctx, userId)
	util.PanicIfError(err)

	if user.ID == 0 {
		return errors.New("user not found")
	}

	return service.userDetailRepository.Delete(ctx, tx, userId)
}

func followToFollowDto(follow []domains.Follow) []dto.FollowDto {
	var followDto []dto.FollowDto
	for _, data := range follow {
		followDto = append(followDto, dto.FollowDto{
			FollowerId:  data.FollowerId,
			FollowingId: data.FollowingId,
			CreatedAt:   data.CreatedAt.Time.String(),
		})
	}

	return followDto
}
