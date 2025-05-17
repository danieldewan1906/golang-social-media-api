package domains

import (
	"context"
	"database/sql"
	"golang-restful-api/dto"
)

type Comment struct {
	ID          int64        `db:"id"`
	UserId      int          `db:"user_id"`
	PostId      int          `db:"post_id"`
	TextComment string       `db:"text_comment"`
	CreatedAt   sql.NullTime `db:"created_at"`
}

type CommentRepository interface {
	FindById(ctx context.Context, id int) (Comment, error)
	FindByPostId(ctx context.Context, postId int) ([]Comment, error)
	FindByUserIdAndPostId(ctx context.Context, userId int, postId int) ([]Comment, error)
	Create(ctx context.Context, tx *sql.Tx, comment Comment) error
	Update(ctx context.Context, tx *sql.Tx, comment Comment) error
	Delete(ctx context.Context, tx *sql.Tx, id int) error
}

type CommentService interface {
	Create(ctx context.Context, req dto.CommentRequestDto) error
	Update(ctx context.Context, id int, req dto.CommentRequestDto) error
	Delete(ctx context.Context, id int) error
}
