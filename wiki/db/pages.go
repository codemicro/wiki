package db

import (
	"database/sql"
	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"time"
)

type Page struct {
	ID        string    `db:"id"`
	Title     string    `db:"title"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Content   string    `db:"content"`
}

func (db *DB) GetAllPages() ([]*Page, error) {
	var o []*Page
	if err := db.pool.Select(&o, `SELECT * FROM "pages"`); err != nil {
		return nil, errors.WithStack(err)
	}
	return o, nil
}

func (db *DB) GetPageByID(id string) (*Page, error) {
	p := new(Page)
	if err := db.pool.Get(p, `SELECT * FROM "pages" WHERE "id" = $1;`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, errors.WithStack(err)
	}
	return p, nil
}

func (db *DB) GetPagesWithTag(tagID string) ([]*Page, error) {
	var o []*Page
	rows, err := db.pool.Queryx(`SELECT "page_id" FROM "page_tag_mapping" WHERE "tag_id" = $1`, tagID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for rows.Next() {
		var pageID string
		if err := rows.Scan(&pageID); err != nil {
			return nil, errors.WithStack(err)
		}
		page, err := db.GetPageByID(pageID)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		o = append(o, page)
	}
	return o, nil
}

func (db *DB) AssignPageToTag(pageID, tagID string) error {
	ctx, cancel := db.newContext()
	defer cancel()

	tx, err := db.pool.BeginTxx(ctx, nil)
	if err != nil {
		return errors.WithStack(err)
	}
	defer smartRollback(tx)

	_, err = tx.Exec(
		`DELETE FROM "page_tag_mapping" WHERE "page_id" = $1; INSERT INTO "page_tag_mapping"("page_id", "tag_id") VALUES ($2, $3)`,
		pageID,
		pageID,
		tagID,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(
		tx.Commit(),
	)
}

func (db *DB) CreatePage(page *Page) error {
	page.CreatedAt = time.Now().UTC()
	page.UpdatedAt = page.CreatedAt

	ctx, cancel := db.newContext()
	defer cancel()

	tx, err := db.pool.BeginTxx(ctx, nil)
	if err != nil {
		return errors.WithStack(err)
	}
	defer smartRollback(tx)

	_, err = tx.NamedExec(
		`INSERT INTO "pages"("id", "title", "created_at", "updated_at", "content") VALUES(:id, :title, :created_at, :updated_at, :content);`,
		page,
	)

	if err != nil {
		if e, ok := err.(sqlite3.Error); ok {
			switch e.ExtendedCode {
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

func (db *DB) UpdatePage(page *Page) error {
	page.UpdatedAt = time.Now().UTC()

	ctx, cancel := db.newContext()
	defer cancel()

	tx, err := db.pool.BeginTxx(ctx, nil)
	if err != nil {
		return errors.WithStack(err)
	}
	defer smartRollback(tx)

	_, err = tx.NamedExec(
		`UPDATE "pages" SET "title" = :title, "updated_at" = :updated_at, "content" = :content WHERE "id" = :id`,
		page,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(
		tx.Commit(),
	)
}
