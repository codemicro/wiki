package db

import (
	"database/sql"
	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type Tag struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

func (db *DB) GetAllTags() ([]*Tag, error) {
	var o []*Tag
	if err := db.pool.Select(&o, `SELECT * FROM "tags";`); err != nil {
		return nil, errors.WithStack(err)
	}
	return o, nil
}

func (db *DB) GetTagByID(id string) (*Tag, error) {
	t := new(Tag)
	if err := db.pool.Get(t, `SELECT * FROM "tags" WHERE "id" = $1;`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, errors.WithStack(err)
	}
	return t, nil
}

func (db *DB) GetTagByName(name string) (*Tag, error) {
	t := new(Tag)
	if err := db.pool.Get(t, `SELECT * FROM "tags" WHERE "name" = $1;`, name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, errors.WithStack(err)
	}
	return t, nil
}

var ErrTagNameExists = errors.New("db: a tag with that name already exists")

func (db *DB) CreateTag(tag *Tag) error {
	ctx, cancel := db.newContext()
	defer cancel()

	tx, err := db.pool.BeginTxx(ctx, nil)
	if err != nil {
		return errors.WithStack(err)
	}
	defer smartRollback(tx)

	_, err = tx.NamedExec(`INSERT INTO "tags"("id", "name") VALUES(:id, :name);`, tag)
	if err != nil {
		if e, ok := err.(sqlite3.Error); ok {
			switch e.ExtendedCode {
			case sqlite3.ErrConstraintUnique:
				return ErrTagNameExists
			case sqlite3.ErrConstraintPrimaryKey:
				return ErrPKAlreadyExists
			}
		}
		return errors.WithStack(err)
	}

	return errors.WithStack(
		tx.Commit(),
	)
}

func (db *DB) UpdateTag(tag *Tag) error {
	ctx, cancel := db.newContext()
	defer cancel()

	tx, err := db.pool.BeginTxx(ctx, nil)
	if err != nil {
		return errors.WithStack(err)
	}
	defer smartRollback(tx)

	_, err = tx.NamedExec(`UPDATE "tags" SET "name" = :name WHERE "id" = :id`, tag)
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(
		tx.Commit(),
	)
}
