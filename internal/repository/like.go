package repository

import (
	"context"
	"database/sql"
	"golang-restful-api/domains"

	"github.com/doug-martin/goqu/v9"
)

type likeRepository struct {
	db *goqu.Database
}

// FindByPostId implements domains.LikeRepository.
func (repository *likeRepository) FindByPostId(ctx context.Context, postId int) (result []domains.Like, err error) {
	data := repository.db.From(goqu.S("golang").Table("likes").As("l")).Where(
		goqu.And(
			goqu.I("l.post_id").Eq(postId),
		),
	)
	err = data.ScanStructsContext(ctx, &result)
	return result, err
}

// FindByUserIdAndPostId implements domains.LikeRepository.
func (repository *likeRepository) FindByUserIdAndPostId(ctx context.Context, userId int, postId int) (result domains.Like, err error) {
	data := repository.db.From(goqu.S("golang").Table("likes").As("l")).Where(
		goqu.And(
			goqu.I("l.user_id").Eq(userId),
			goqu.I("l.post_id").Eq(postId),
		),
	).Limit(1)
	_, err = data.ScanStructContext(ctx, &result)
	return result, err
}

// Create implements domains.LikeRepository.
func (repository *likeRepository) Create(ctx context.Context, tx *sql.Tx, like domains.Like) error {
	sql, args, _ := repository.db.Insert(goqu.S("golang").Table("likes")).Rows(
		goqu.Record{
			"user_id":    like.UserId,
			"post_id":    like.PostId,
			"created_at": like.CreatedAt,
		},
	).ToSQL()
	_, err := tx.ExecContext(ctx, sql, args...)
	return err
}

// Delete implements domains.LikeRepository.
func (repository *likeRepository) Delete(ctx context.Context, tx *sql.Tx, userId int, postId int) error {
	sql, args, _ := repository.db.Delete(goqu.S("golang").Table("likes")).Where(
		goqu.And(
			goqu.C("user_id").Eq(userId),
			goqu.C("post_id").Eq(postId),
		),
	).ToSQL()
	_, err := tx.ExecContext(ctx, sql, args...)
	return err
}

func NewLikeRepository(db *sql.DB) domains.LikeRepository {
	return &likeRepository{
		db: goqu.New("default", db),
	}
}
