package dbuser

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/perfectgentlemande/go-postgresql-crud-example/internal/service"
)

const usersTable = "users"

type User struct {
	ID        string    `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (u *User) ToService() service.User {
	return service.User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
func userFromService(u *service.User) User {
	return User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

type Users []User

func (us Users) ToService() []service.User {
	res := make([]service.User, 0, len(us))

	for i := range us {
		res = append(res, us[i].ToService())
	}

	return res
}

func (d *Database) CreateUser(ctx context.Context, user *service.User) (string, error) {
	dbUser := userFromService(user)

	sqlText, bound, err := squirrel.Insert(usersTable).Columns(
		"id",
		"username",
		"email",
		"created_at",
		"updated_at").Values(
		dbUser.ID,
		dbUser.Username,
		dbUser.Email,
		dbUser.CreatedAt,
		dbUser.UpdatedAt).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return "", fmt.Errorf("cannot build sql: %w", err)
	}

	_, err = d.db.ExecContext(ctx, sqlText, bound...)
	if err != nil {
		return "", fmt.Errorf("cannot insert user: %w", err)
	}

	return user.ID, nil
}

func withListUsersParams(query squirrel.SelectBuilder, params *service.ListUsersParams) squirrel.SelectBuilder {
	if params.Limit != nil {
		query = query.Limit(*params.Limit)
	}
	if params.Offset != nil {
		query = query.Offset(*params.Offset)
	}

	return query
}

func (d *Database) ListUsers(ctx context.Context, params *service.ListUsersParams) ([]service.User, error) {
	var res Users

	query := squirrel.Select(
		"id",
		"username",
		"email",
		"created_at",
		"updated_at").
		From(usersTable).
		PlaceholderFormat(squirrel.Dollar)
	query = withListUsersParams(query, params)

	sqlText, bound, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("cannot build SQL: %w", err)
	}

	if err = d.db.SelectContext(ctx, &res, sqlText, bound...); err != nil {
		return nil, fmt.Errorf("cannot select users: %w", err)
	}

	return res.ToService(), nil
}

func (d *Database) GetUserByID(ctx context.Context, id string) (service.User, error) {
	res := User{}

	query := squirrel.Select(
		"id",
		"username",
		"email",
		"created_at",
		"updated_at").
		From(usersTable).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)

	sqlText, bound, err := query.ToSql()
	if err != nil {
		return service.User{}, fmt.Errorf("cannot build SQL: %w", err)
	}

	if err = d.db.GetContext(ctx, &res, sqlText, bound...); err != nil {
		return service.User{}, fmt.Errorf("cannot get user: %w", err)
	}

	return res.ToService(), nil
}

func (d *Database) UpdateUserByID(ctx context.Context, id string, user *service.User) (service.User, error) {
	dbUser := userFromService(user)

	sqlText, bound, err := squirrel.Update(usersTable).SetMap(map[string]interface{}{
		"username":   dbUser.Username,
		"email":      dbUser.Email,
		"updated_at": dbUser.UpdatedAt,
	}).Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return service.User{}, fmt.Errorf("cannot build SQL: %w", err)
	}

	_, err = d.db.ExecContext(ctx, sqlText, bound...)
	if err != nil {
		return service.User{}, fmt.Errorf("cannot update user: %w", err)
	}

	return *user, nil
}

func (d *Database) DeleteUserByID(ctx context.Context, id string) error {
	sqlText, bound, err := squirrel.Delete(usersTable).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("cannot build SQL: %w", err)
	}

	_, err = d.db.ExecContext(ctx, sqlText, bound...)
	if err != nil {
		return fmt.Errorf("cannot delete user: %w", err)
	}

	return nil
}
