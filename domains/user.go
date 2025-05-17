package domains

import (
	"context"
	"database/sql"
)

type User struct {
	ID        int64          `db:"id"`
	Email     string         `db:"email"`
	Username  string         `db:"username"`
	Password  string         `db:"password"`
	Token     sql.NullString `db:"token"`
	IsActive  bool           `db:"is_active"`
	CreatedAt sql.NullTime   `db:"created_at"`
	Role      string         `db:"role"`
}

type UserRepository interface {
	FindAll(ctx context.Context) ([]User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByEmailUsername(ctx context.Context, email string, username string) (User, error)
	Save(ctx context.Context, tx *sql.Tx, usr *User) (User, error)
	Update(ctx context.Context, tx *sql.Tx, usr *User) error
	Delete(ctx context.Context, tx *sql.Tx, id int) error
}
