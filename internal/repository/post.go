package repository

import (
	"context"
	"database/sql"
	"golang-restful-api/domains"

	"github.com/doug-martin/goqu/v9"
)

type postRepository struct {
	db *goqu.Database
}

// FindByUserId implements domains.PostRepository.
func (repository *postRepository) FindByUserId(ctx context.Context, userId int) (result []domains.Post, err error) {
	dataset := repository.db.From(goqu.S("golang").Table("posts").As("p")).Where(
		goqu.C("user_id").Eq(userId),
	)
	err = dataset.ScanStructsContext(ctx, &result)
	return result, err
}

// FindByUserIdAndId implements domains.PostRepository.
func (repository *postRepository) FindByUserIdAndId(ctx context.Context, userId int, id int) (result domains.Post, err error) {
	dataset := repository.db.From(goqu.S("golang").Table("posts").As("p")).Where(
		goqu.And(
			goqu.I("p.id").Eq(id),
			goqu.I("p.user_id").Eq(userId),
		),
	)
	_, err = dataset.ScanStructContext(ctx, &result)
	return result, err
}

// ArchievePost implements domains.PostRepository.
func (repository *postRepository) ArchievePost(ctx context.Context, tx *sql.Tx, id int) error {
	sql, args, _ := repository.db.Update(goqu.S("golang").Table("posts")).Set(goqu.Record{
		"is_active": false,
	}).Where(
		goqu.C("id").Eq(id),
	).ToSQL()
	_, err := tx.ExecContext(ctx, sql, args...)
	return err
}

// Create implements domains.PostRepository.
func (repository *postRepository) Create(ctx context.Context, tx *sql.Tx, post domains.Post) error {
	sql, args, _ := repository.db.Insert(goqu.S("golang").Table("posts")).Rows(
		goqu.Record{
			"user_id":    post.UserId,
			"content":    post.Content,
			"image_url":  post.ImageUrl,
			"is_active":  post.IsActive,
			"created_at": post.CreatedAt,
		},
	).ToSQL()
	_, err := tx.ExecContext(ctx, sql, args...)
	return err
}

// Delete implements domains.PostRepository.
func (repository *postRepository) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	sql, args, _ := repository.db.Delete(goqu.S("golang").Table("posts")).Where(
		goqu.C("id").Eq(id),
	).ToSQL()
	_, err := tx.ExecContext(ctx, sql, args...)
	return err
}

// FindAll implements domains.PostRepository.
func (repository *postRepository) FindAll(ctx context.Context) (result []domains.Post, err error) {
	dataset := repository.db.From(goqu.S("golang").Table("posts").As("p")).Where(goqu.I("p.is_active").IsTrue())
	err = dataset.ScanStructsContext(ctx, &result)
	return result, err
}

// FindById implements domains.PostRepository.
func (repository *postRepository) FindById(ctx context.Context, id int) (result domains.Post, err error) {
	dataset := repository.db.From(goqu.S("golang").Table("posts").As("p")).Where(
		goqu.And(
			goqu.I("p.is_active").IsTrue(),
			goqu.I("p.id").Eq(id),
		),
	)
	_, err = dataset.ScanStructContext(ctx, &result)
	return result, err
}

// Update implements domains.PostRepository.
func (repository *postRepository) Update(ctx context.Context, tx *sql.Tx, post domains.Post) error {
	sql, args, _ := repository.db.Update(goqu.S("golang").Table("posts").As("p")).Where(
		goqu.And(
			goqu.I("p.id").Eq(post.ID),
			goqu.I("p.user_id").Eq(post.UserId),
		),
	).Set(
		goqu.Record{
			"content":    post.Content,
			"image_url":  post.ImageUrl,
			"is_active":  post.IsActive,
			"updated_at": post.UpdatedAt,
		},
	).ToSQL()
	_, err := tx.ExecContext(ctx, sql, args...)
	return err
}

func NewPostRepository(conn *sql.DB) domains.PostRepository {
	return &postRepository{
		db: goqu.New("default", conn),
	}
}
