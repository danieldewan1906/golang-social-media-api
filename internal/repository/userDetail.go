package repository

import (
	"context"
	"database/sql"
	"golang-restful-api/domains"
	"golang-restful-api/dto"
	"strings"

	"github.com/doug-martin/goqu/v9"
)

type userDetailRepository struct {
	db *goqu.Database
}

// Delete implements domain.userDetailRepository.
func (cr *userDetailRepository) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	sql, args, _ := cr.db.Update("users").
		Where(goqu.C("id").Eq(id)).
		Set(goqu.Record{
			"is_active": false,
		}).ToSQL()
	_, err := tx.ExecContext(ctx, sql, args...)
	return err
}

// FindAll implements domain.userDetailRepository.
func (cr userDetailRepository) FindAll(ctx context.Context, req dto.UserDetailRequestDto) (result []domains.UserDetail, err error) {
	var where goqu.Expression

	if req.ID == 0 && req.Name == "" {
		where = goqu.L("1 = 1")
	}

	if req.ID != 0 {
		where = goqu.I("u.id").Eq(req.ID)
	}

	if req.Name != "" {
		where = goqu.L("CONCAT(?, ' ', ?)", goqu.I("ud.first_name"), goqu.I("ud.last_name")).ILike("%" + strings.ToLower(req.Name) + "%")
	}

	dataset := cr.db.
		Select(
			goqu.I("u.id"),
			goqu.I("ud.first_name"),
			goqu.I("ud.last_name"),
			goqu.I("ud.address"),
			goqu.I("ud.created_at"),
			goqu.I("ud.updated_at"),
		).
		From(goqu.S("golang").Table("users").As("u")).
		Join(
			goqu.S("golang").Table("user_details").As("ud"),
			goqu.On(goqu.Ex{
				"u.id": goqu.I("ud.user_id")})).
		Where(
			goqu.And(
				goqu.I("u.is_active").IsTrue(),
				goqu.Or(
					where,
				),
			),
		)
	err = dataset.ScanStructsContext(ctx, &result)
	return
}

// FindById implements domain.userDetailRepository.
func (cr *userDetailRepository) FindByUserId(ctx context.Context, userId int) (result domains.UserDetail, err error) {
	dataset := cr.db.From("user_details").
		Where(goqu.C("user_id").Eq(userId))
	_, err = dataset.ScanStructContext(ctx, &result)
	return
}

// Save implements domain.userDetailRepository.
func (cr *userDetailRepository) Save(ctx context.Context, tx *sql.Tx, c *domains.UserDetail) error {
	sql, args, _ := cr.db.Insert("user_details").Rows(goqu.Record{
		"user_id":    c.UserId,
		"first_name": c.FirstName,
		"last_name":  c.LastName,
		"address":    c.Address,
		"created_at": c.CreatedAt,
		"updated_at": c.UpdatedAt,
	}).ToSQL()
	_, err := tx.ExecContext(ctx, sql, args...)
	return err
}

// Update implements domain.userDetailRepository.
func (cr *userDetailRepository) Update(ctx context.Context, tx *sql.Tx, c *domains.UserDetail) error {
	sql, args, _ := cr.db.Update("user_details").Where(
		goqu.And(
			goqu.C("id").Eq(c.ID),
			goqu.C("user_id").Eq(c.UserId),
		),
	).Set(goqu.Record{
		"first_name": c.FirstName,
		"last_name":  c.LastName,
		"address":    c.Address,
		"updated_at": c.UpdatedAt,
	}).ToSQL()
	_, err := tx.ExecContext(ctx, sql, args...)
	return err
}

func NewUserDetail(con *sql.DB) domains.UserDetailRepository {
	return &userDetailRepository{
		db: goqu.New("default", con),
	}
}
