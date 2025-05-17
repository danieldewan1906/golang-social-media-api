package repository

import (
	"context"
	"database/sql"
	"golang-restful-api/domains"

	"github.com/doug-martin/goqu/v9"
)

type userRepository struct {
	db *goqu.Database
}

func NewUser(conn *sql.DB) domains.UserRepository {
	return &userRepository{
		db: goqu.New("default", conn),
	}
}

// FindByEmail implements domains.UserRepository.
func (repository *userRepository) FindByEmail(ctx context.Context, email string) (usr domains.User, err error) {
	dataset := repository.db.From("users").Where(goqu.C("email").Eq(email))
	_, err = dataset.ScanStructContext(ctx, &usr)
	return
}

// FindByEmailUsername implements domains.UserRepository.
func (repository *userRepository) FindByEmailUsername(ctx context.Context, email string, username string) (usr domains.User, err error) {
	dataset := repository.db.From("users").Where(
		goqu.And(
			goqu.C("is_active").IsTrue(),
			goqu.Or(
				goqu.C("email").Eq(email),
				goqu.C("username").Eq(username),
			),
		),
	)
	_, err = dataset.ScanStructContext(ctx, &usr)
	return
}

// Save implements domains.UserRepository.
func (repository *userRepository) Save(ctx context.Context, tx *sql.Tx, usr *domains.User) (domains.User, error) {
	sql, args, _ := repository.db.Insert("users").Rows(goqu.Record{
		"username":   usr.Username,
		"email":      usr.Email,
		"password":   usr.Password,
		"token":      usr.Token,
		"is_active":  usr.IsActive,
		"created_at": usr.CreatedAt,
		"role":       usr.Role,
	}).Returning("id").ToSQL()
	result := tx.QueryRowContext(ctx, sql, args...)
	return *usr, result.Scan(&usr.ID)
}

// Update implements domains.UserRepository.
func (repository *userRepository) Update(ctx context.Context, tx *sql.Tx, usr *domains.User) error {
	sql, args, _ := repository.db.Update("users").Where(goqu.C("id").Eq(usr.ID)).Set(goqu.Record{
		"username":  usr.Username,
		"email":     usr.Email,
		"password":  usr.Password,
		"token":     usr.Token,
		"is_active": usr.IsActive,
		"role":      usr.Role,
	}).ToSQL()
	_, err := tx.ExecContext(ctx, sql, args...)
	return err
}

// Delete implements domains.UserRepository.
func (repository *userRepository) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	sql, _, _ := repository.db.Update("users").Where(goqu.C("id").Eq(id)).Set(goqu.Record{
		"is_active": false,
	}).ToSQL()
	_, err := tx.ExecContext(ctx, sql, []interface{}{}...)
	return err
}

// FindAll implements domains.UserRepository.
func (repository *userRepository) FindAll(ctx context.Context) (result []domains.User, err error) {
	data := repository.db.From("users").Where(goqu.C("is_active").Eq(true))
	err = data.ScanStructsContext(ctx, &result)
	return result, err
}
