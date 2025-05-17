package repository

import (
	"context"
	"database/sql"
	"golang-restful-api/domains"

	"github.com/doug-martin/goqu/v9"
)

type userImageRepository struct {
	db *goqu.Database
}

// Delete implements domains.UserImagesRepository.
func (repository *userImageRepository) Delete(ctx context.Context, tx *sql.Tx, userId int) error {
	sql, args, _ := repository.db.Delete("user_images").Where(goqu.C("user_id").Eq(userId)).ToSQL()
	_, err := tx.ExecContext(ctx, sql, args...)
	return err
}

// FindByUserId implements domains.UserImagesRepository.
func (repository *userImageRepository) FindByUserId(ctx context.Context, userId int) (result domains.UserImages, err error) {
	dataset := repository.db.From("user_images").Where(goqu.C("user_id").Eq(userId))
	_, err = dataset.ScanStructContext(ctx, &result)
	return result, err
}

// Save implements domains.UserImagesRepository.
func (repository *userImageRepository) Save(ctx context.Context, tx *sql.Tx, userImages domains.UserImages) error {
	sql, args, _ := repository.db.Insert("user_images").Rows(goqu.Record{
		"user_id":    userImages.UserId,
		"image_url":  userImages.ImageUrl,
		"created_at": userImages.CreatedAt,
		"updated_at": userImages.UpdatedAt,
		"extension":  userImages.Extension,
	}).ToSQL()
	_, err := tx.ExecContext(ctx, sql, args...)
	return err
}

// Update implements domains.UserImagesRepository.
func (repository *userImageRepository) Update(ctx context.Context, tx *sql.Tx, userImages domains.UserImages) error {
	sql, args, _ := repository.db.Update("user_images").Where(
		goqu.And(
			goqu.C("id").Eq(userImages.ID),
			goqu.C("user_id").Eq(userImages.UserId),
		),
	).Set(goqu.Record{
		"image_url":  userImages.ImageUrl,
		"updated_at": userImages.UpdatedAt,
		"extension":  userImages.Extension,
	}).ToSQL()
	_, err := tx.ExecContext(ctx, sql, args...)
	return err
}

func NewUserImages(con *sql.DB) domains.UserImagesRepository {
	return &userImageRepository{
		db: goqu.New("default", con),
	}
}
