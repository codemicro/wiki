package db

import (
	"context"
	"database/sql"
	"github.com/codemicro/wiki/wiki/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"math"
	"net"
	"time"
)

var ErrNotFound = errors.New("db: record not found")

type DB struct {
	pool           *sqlx.DB
	ContextTimeout time.Duration
}

const maxConnectionAttempts = 4

func New() (*DB, error) {
	dsn := config.Database.Filename
	log.Info().Msg("connecting to database")
	db, err := sqlx.Open("sqlite3", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "could not open SQL connection")
	}

	rtn := &DB{
		pool:           db,
		ContextTimeout: time.Second,
	}

	for i := 1; i <= maxConnectionAttempts; i += 1 {
		logger := log.With().Int("attempt", i).Int("maxAttempts", maxConnectionAttempts).Logger()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		err := rtn.pool.PingContext(ctx)

		if err == nil {
			cancel()
			break
		}

		if e, ok := err.(*net.OpError); ((ok && e.Op == "dial") || errors.Is(err, context.DeadlineExceeded)) && i != maxConnectionAttempts {
			cancel()

			retryIn := int(math.Pow(math.E, float64(i)))
			logger.Warn().Err(err).Msgf("could not connect to database - retrying in %d seconds", retryIn)
			time.Sleep(time.Second * time.Duration(retryIn))

			continue
		}

		cancel()
		return nil, errors.Wrapf(err, "could not ping database after %d attempts", i)
	}

	return rtn, nil
}

func (db *DB) newContext() (context.Context, func()) {
	return context.WithTimeout(context.Background(), db.ContextTimeout)
}

func smartRollback(tx *sqlx.Tx) {
	err := tx.Rollback()
	if err != nil && !errors.Is(err, sql.ErrTxDone) {
		log.Warn().Stack().Err(errors.WithStack(err)).Str("location", "smartRollback").Msg("failed to rollback transaction")
	}
}
