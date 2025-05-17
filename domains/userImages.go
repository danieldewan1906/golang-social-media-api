package domains

import (
	"context"
	"database/sql"
	"golang-restful-api/dto"
)

type UserImages struct {
	ID        int64          `db:"id"`
	UserId    int            `db:"user_id"`
	ImageUrl  sql.NullString `db:"image_url"`
	CreatedAt sql.NullTime   `db:"created_at"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
	Extension sql.NullString `db:"extension"`
}

type UserImagesRepository interface {
	FindByUserId(ctx context.Context, userId int) (UserImages, error)
	Save(ctx context.Context, tx *sql.Tx, userImages UserImages) error
	Update(ctx context.Context, tx *sql.Tx, userImages UserImages) error
	Delete(ctx context.Context, tx *sql.Tx, userId int) error
}

type UserImagesService interface {
	FindByUserId(ctx context.Context, userId int) (dto.UserImagesDto, error)
	Save(ctx context.Context, req dto.UserImageRequestDto, userId int) error
	Update(ctx context.Context, req dto.UserImageRequestDto, userId int) error
	Delete(ctx context.Context, userId int) error
}
