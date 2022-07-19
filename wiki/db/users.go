package db

import (
	"database/sql"
	"github.com/pkg/errors"
)

type User struct {
	ID         string         `db:"id"`
	ExternalID string         `db:"external_id"`
	Name       sql.NullString `db:"name"`
	Email      string         `db:"email"`
}

func (db *DB) GetUserByID(id string) (*User, error) {
	u := new(User)
	err := db.pool.Get(u, `SELECT * FROM "users" WHERE "id" = $1`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, errors.WithStack(err)
	}
	return u, nil
}

func (db *DB) GetUserByExternalID(externalID string) (*User, error) {
	u := new(User)
	err := db.pool.Get(u, `SELECT * FROM "users" WHERE "external_id" = $1`, externalID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, errors.WithStack(err)
	}
	return u, nil
}

func (db *DB) CreateUser(u *User) error {
	ctx, cancel := db.newContext()
	defer cancel()

	tx, err := db.pool.BeginTxx(ctx, nil)
	if err != nil {
		return errors.WithStack(err)
	}
	defer smartRollback(tx)

	_, err = tx.NamedExec(
		`INSERT INTO "users"("id", "external_id", "name", "email") VALUES (:id, :external_id, :name, :email)`,
		u,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(
		tx.Commit(),
	)
}
