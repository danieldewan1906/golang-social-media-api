package domains

import (
	"context"
	"database/sql"
	"golang-restful-api/dto"
)

type Post struct {
	ID        int64          `db:"id"`
	UserId    int            `db:"user_id"`
	Content   sql.NullString `db:"content"`
	ImageUrl  sql.NullString `db:"image_url"`
	IsActive  bool           `db:"is_active"`
	CreatedAt sql.NullTime   `db:"created_at"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
}

type PostRepository interface {
	FindAll(ctx context.Context) ([]Post, error)
	FindById(ctx context.Context, id int) (Post, error)
	FindByUserId(ctx context.Context, userId int) ([]Post, error)
	FindByUserIdAndId(ctx context.Context, userId int, id int) (Post, error)
	Create(ctx context.Context, tx *sql.Tx, post Post) error
	Update(ctx context.Context, tx *sql.Tx, post Post) error
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	ArchievePost(ctx context.Context, tx *sql.Tx, id int) error
}

type PostService interface {
	FindAll(ctx context.Context) ([]dto.PostDto, error)
	FindById(ctx context.Context, id int) (dto.PostDto, error)
	FindByUserId(ctx context.Context, userId int) ([]dto.PostDto, error)
	Create(ctx context.Context, req dto.PostRequestDto) error
	Update(ctx context.Context, id int, req dto.PostRequestDto) error
	Delete(ctx context.Context, id int) error
	ArchievePost(ctx context.Context, id int) error
}
