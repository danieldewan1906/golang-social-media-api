package repository

import (
	"context"
	"database/sql"
	"golang-restful-api/domains"
	"log"

	"github.com/doug-martin/goqu/v9"
)

type followRepository struct {
	db *goqu.Database
}

// FindByFollowerId implements domains.FollowRepository.
func (repository *followRepository) FindByFollowerId(ctx context.Context, followerId int) (result []domains.Follow, err error) {
	data := repository.db.From(goqu.S("golang").Table("follows")).Where(
		goqu.C("follower_id").Eq(followerId),
	)

	log.Println(data.ToSQL())

	err = data.ScanStructsContext(ctx, &result)
	return result, err
}

// FindByFollowingId implements domains.FollowRepository.
func (repository *followRepository) FindByFollowingId(ctx context.Context, followingId int) (result []domains.Follow, err error) {
	data := repository.db.From(goqu.S("golang").Table("follows")).Where(
		goqu.C("following_id").Eq(followingId),
	)

	log.Println(data.ToSQL())

	err = data.ScanStructsContext(ctx, &result)
	return result, err
}

// FindByFollowingAndFollowerId implements domains.FollowRepository.
func (repository *followRepository) FindByFollowingAndFollowerId(ctx context.Context, followerId int, followingId int) (result domains.Follow, err error) {
	data := repository.db.From(goqu.S("golang").Table("follows")).Where(
		goqu.And(
			goqu.C("follower_id").Eq(followerId),
			goqu.C("following_id").Eq(followingId),
		),
	)

	log.Println(data.ToSQL())

	_, err = data.ScanStructContext(ctx, &result)
	return result, err
}

// Create implements domains.FollowRepository.
func (repository *followRepository) Create(ctx context.Context, tx *sql.Tx, follow domains.Follow) error {
	sql, args, _ := repository.db.Insert(goqu.S("golang").Table("follows")).Rows(
		goqu.Record{
			"follower_id":  follow.FollowerId,
			"following_id": follow.FollowingId,
			"created_at":   follow.CreatedAt,
		},
	).ToSQL()

	log.Println(sql)

	_, err := tx.ExecContext(ctx, sql, args...)
	return err
}

// Delete implements domains.FollowRepository.
func (repository *followRepository) Delete(ctx context.Context, tx *sql.Tx, follow domains.Follow) error {
	sql, args, _ := repository.db.Delete(goqu.S("golang").Table("follows")).Where(
		goqu.And(
			goqu.C("follower_id").Eq(follow.FollowerId),
			goqu.C("following_id").Eq(follow.FollowingId),
		),
	).ToSQL()

	log.Println(sql)

	_, err := tx.ExecContext(ctx, sql, args...)
	return err
}

func NewFollowRepository(conn *sql.DB) domains.FollowRepository {
	return &followRepository{
		db: goqu.New("default", conn),
	}
}
