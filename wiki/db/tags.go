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

func (db *DB) GetTagsByPageID(pageID string) ([]*Tag, error) {
	var o []*Tag
	rows, err := db.pool.Queryx(`SELECT "tag_id" FROM "page_tag_mapping" WHERE "page_id" = $1`, pageID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for rows.Next() {
		var tagID string
		if err := rows.Scan(&tagID); err != nil {
			return nil, errors.WithStack(err)
		}
		tag, err := db.GetTagByID(tagID)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		o = append(o, tag)
	}
	return o, nil
}

func (db *DB) GetTagFrequency(tag *Tag) (int, error) {
	var f int
	err := db.pool.Get(&f, `SELECT COUNT(*) FROM "page_tag_mapping" WHERE "tag_id"=$1`, tag.ID)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return f, nil
}

func (db *DB) GetTagFrequencies(tags []*Tag) (map[*Tag]int, error) {
	freqs := make(map[*Tag]int)
	for _, tag := range tags {
		f, err := db.GetTagFrequency(tag)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		freqs[tag] = f
	}
	return freqs, nil
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
