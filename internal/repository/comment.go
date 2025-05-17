package repository

import (
	"context"
	"database/sql"
	"golang-restful-api/domains"
	"log"

	"github.com/doug-martin/goqu/v9"
)

type commentRepository struct {
	DB *goqu.Database
}

// FindByPostId implements domains.CommentRepository.
func (repository *commentRepository) FindByPostId(ctx context.Context, postId int) (result []domains.Comment, err error) {
	data := repository.DB.From(goqu.S("golang").Table("comments")).Where(goqu.C("post_id").Eq(postId))
	err = data.ScanStructsContext(ctx, &result)
	return result, err
}

// FindById implements domains.CommentRepository.
func (repository *commentRepository) FindById(ctx context.Context, id int) (result domains.Comment, err error) {
	data := repository.DB.From(goqu.S("golang").Table("comments")).Where(goqu.C("id").Eq(id))
	_, err = data.ScanStructContext(ctx, &result)
	return result, err
}

// FindByUserIdAndPostId implements domains.CommentRepository.
func (repository *commentRepository) FindByUserIdAndPostId(ctx context.Context, userId int, postId int) (result []domains.Comment, err error) {
	data := repository.DB.From(goqu.S("golang").Table("comments")).Where(
		goqu.And(
			goqu.C("user_id").Eq(userId),
			goqu.C("post_id").Eq(postId),
		),
	)
	err = data.ScanStructsContext(ctx, &result)
	return result, err
}

// Create implements domains.CommentRepository.
func (repository *commentRepository) Create(ctx context.Context, tx *sql.Tx, comment domains.Comment) error {
	sql, args, _ := repository.DB.Insert(goqu.S("golang").Table("comments")).Rows(
		goqu.Record{
			"user_id":      comment.UserId,
			"post_id":      comment.PostId,
			"text_comment": comment.TextComment,
			"created_at":   comment.CreatedAt,
		},
	).ToSQL()

	log.Println(sql)

	_, err := tx.ExecContext(ctx, sql, args...)
	return err
}

// Delete implements domains.CommentRepository.
func (repository *commentRepository) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	sql, args, _ := repository.DB.Delete(goqu.S("golang").Table("comments")).Where(
		goqu.C("id").Eq(id),
	).ToSQL()

	log.Println(sql)

	_, err := tx.ExecContext(ctx, sql, args...)
	return err
}

// Update implements domains.CommentRepository.
func (repository *commentRepository) Update(ctx context.Context, tx *sql.Tx, comment domains.Comment) error {
	sql, args, _ := repository.DB.Update(goqu.S("golang").Table("comments")).Where(goqu.C("id").Eq(comment.ID)).Set(
		goqu.Record{
			"text_comment": comment.TextComment,
		},
	).ToSQL()

	log.Println(sql)

	_, err := tx.ExecContext(ctx, sql, args...)
	return err
}

func NewCommentRepository(conn *sql.DB) domains.CommentRepository {
	return &commentRepository{
		DB: goqu.New("default", conn),
	}
}
