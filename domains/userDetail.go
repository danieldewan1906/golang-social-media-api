package domains

import (
	"context"
	"database/sql"
	"golang-restful-api/dto"
)

type UserDetail struct {
	ID        int64          `db:"id"`
	UserId    int            `db:"user_id"`
	FirstName string         `db:"first_name"`
	LastName  sql.NullString `db:"last_name"`
	Address   sql.NullString `db:"address"`
	CreatedAt sql.NullTime   `db:"created_at"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
}

type UserDetailRepository interface {
	FindAll(ctx context.Context, req dto.UserDetailRequestDto) ([]UserDetail, error)
	FindByUserId(ctx context.Context, userId int) (UserDetail, error)
	Save(ctx context.Context, tx *sql.Tx, c *UserDetail) error
	Update(ctx context.Context, tx *sql.Tx, c *UserDetail) error
	Delete(ctx context.Context, tx *sql.Tx, id int) error
}

type UserDetailService interface {
	FindAll(ctx context.Context, req dto.UserDetailRequestDto) ([]dto.UserDetailDto, error)
	FindByUserId(ctx context.Context, userId int) (dto.UserDetailDto, error)
	Update(ctx context.Context, req dto.UpdateUserRequestDto) error
	InActiveUser(ctx context.Context, userId int) error
}
