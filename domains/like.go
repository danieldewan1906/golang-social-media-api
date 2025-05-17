package domains

import (
	"context"
	"database/sql"
	"golang-restful-api/dto"
)

type Like struct {
	ID        int64        `db:"id"`
	UserId    int          `db:"user_id"`
	PostId    int          `db:"post_id"`
	CreatedAt sql.NullTime `db:"created_at"`
}

type LikeRepository interface {
	FindByPostId(ctx context.Context, postId int) ([]Like, error)
	FindByUserIdAndPostId(ctx context.Context, userId int, postId int) (Like, error)
	Create(ctx context.Context, tx *sql.Tx, like Like) error
	Delete(ctx context.Context, tx *sql.Tx, userId int, postId int) error
}

type LikeService interface {
	Create(ctx context.Context, req dto.LikeRequestDto) error
	Delete(ctx context.Context, req dto.LikeRequestDto) error
}
