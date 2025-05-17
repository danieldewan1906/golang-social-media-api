package domains

import (
	"context"
	"database/sql"
	"golang-restful-api/dto"
)

type Follow struct {
	FollowerId  int          `db:"follower_id"`
	FollowingId int          `db:"following_id"`
	CreatedAt   sql.NullTime `db:"created_at"`
}

type FollowRepository interface {
	FindByFollowerId(ctx context.Context, followerId int) ([]Follow, error)
	FindByFollowingId(ctx context.Context, followingId int) ([]Follow, error)
	FindByFollowingAndFollowerId(ctx context.Context, followerId int, followingId int) (Follow, error)
	Create(ctx context.Context, tx *sql.Tx, follow Follow) error
	Delete(ctx context.Context, tx *sql.Tx, follow Follow) error
}

type FollowService interface {
	Create(ctx context.Context, req dto.FollowRequestDto) error
	Delete(ctx context.Context, req dto.FollowRequestDto) error
}
