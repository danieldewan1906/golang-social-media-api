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

type userImageService struct {
	DB                  *sql.DB
	userImageRepository domains.UserImagesRepository
}

// Delete implements domains.UserImagesService.
func (service *userImageService) Delete(ctx context.Context, userId int) error {
	tx, err := service.DB.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	userImage, err := service.userImageRepository.FindByUserId(ctx, userId)
	util.PanicIfError(err)
	if userImage.ID == 0 || userImage.UserId != userId {
		util.PanicIfError(errors.New("images not found"))
	}

	return service.userImageRepository.Delete(ctx, tx, userImage.UserId)
}

// FindByUserId implements domains.UserImagesService.
func (service *userImageService) FindByUserId(ctx context.Context, userId int) (dto.UserImagesDto, error) {
	usrImage, err := service.userImageRepository.FindByUserId(ctx, userId)
	if err != nil {
		return dto.UserImagesDto{}, err
	}

	if usrImage.ID == 0 {
		return dto.UserImagesDto{}, errors.New("image not found")
	}

	return dto.UserImagesDto{
		ID:        usrImage.ID,
		UserId:    usrImage.UserId,
		ImageUrl:  usrImage.ImageUrl.String,
		CreatedAt: usrImage.CreatedAt.Time.String(),
		UpdatedAt: usrImage.UpdatedAt.Time.String(),
		Extension: usrImage.Extension.String,
		BaseUrl:   "/Users/daniel/Documents/learn/golang/uploads",
	}, nil
}

// Save implements domains.UserImagesService.
func (service *userImageService) Save(ctx context.Context, req dto.UserImageRequestDto, userId int) error {
	tx, err := service.DB.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	usrImage, err := service.userImageRepository.FindByUserId(ctx, userId)
	util.PanicIfError(err)

	if usrImage.ID != 0 {
		util.PanicIfError(errors.New("images already set"))
	}

	imageUrl := "/" + req.Filename
	reqUserImg := domains.UserImages{
		UserId:    userId,
		ImageUrl:  sql.NullString{Valid: true, String: imageUrl},
		CreatedAt: sql.NullTime{Valid: true, Time: time.Now()},
		Extension: sql.NullString{Valid: true, String: req.Extension},
	}

	return service.userImageRepository.Save(ctx, tx, reqUserImg)
}

// Update implements domains.UserImagesService.
func (service *userImageService) Update(ctx context.Context, req dto.UserImageRequestDto, userId int) error {
	tx, err := service.DB.Begin()
	util.PanicIfError(err)
	defer util.CommitOrRollback(tx)

	userImage, err := service.userImageRepository.FindByUserId(ctx, userId)
	util.PanicIfError(err)
	if userImage.ID == 0 || userImage.UserId != userId {
		util.PanicIfError(errors.New("images not found"))
	}

	imageUrl := "/" + req.Filename
	reqUserImg := domains.UserImages{
		ID:        userImage.ID,
		UserId:    userImage.UserId,
		ImageUrl:  sql.NullString{Valid: true, String: imageUrl},
		UpdatedAt: sql.NullTime{Valid: true, Time: time.Now()},
		Extension: sql.NullString{Valid: true, String: req.Extension},
	}

	return service.userImageRepository.Update(ctx, tx, reqUserImg)
}

func NewUserImageService(db *sql.DB, userImageRepository domains.UserImagesRepository) domains.UserImagesService {
	return &userImageService{
		DB:                  db,
		userImageRepository: userImageRepository,
	}
}
