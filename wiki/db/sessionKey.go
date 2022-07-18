package db

import (
	"database/sql"
	"github.com/pkg/errors"
)

func (db *DB) StoreSessionKey(key string) error {
	ctx, cancel := db.newContext()
	defer cancel()

	tx, err := db.pool.BeginTxx(ctx, nil)
	if err != nil {
		return errors.WithStack(err)
	}
	defer smartRollback(tx)

	if _, err := tx.Exec(`DELETE FROM "session_key"`); err != nil {
		return errors.WithStack(err)
	}

	if _, err := tx.Exec(`INSERT INTO "session_key"("key") VALUES ($1)`, key); err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(
		tx.Commit(),
	)
}

func (db *DB) GetSessionKey() (string, error) {
	var key string
	if err := db.pool.QueryRowx(`SELECT * FROM "session_key"`).Scan(&key); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNotFound
		}
	}
	return key, nil
}
