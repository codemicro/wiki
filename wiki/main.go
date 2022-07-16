package main

import (
	"fmt"
	"github.com/codemicro/wiki/wiki/config"
	"github.com/codemicro/wiki/wiki/db"
	"github.com/codemicro/wiki/wiki/endpoints"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strconv"
)

func run() error {
	database, err := db.New()
	if err != nil {
		return errors.WithStack(err)
	}

	if err := database.Migrate(); err != nil {
		return errors.Wrap(err, "failed migration")
	}

	e := endpoints.New(nil)
	app := e.SetupApp()

	serveAddr := config.HTTP.Host + ":" + strconv.Itoa(config.HTTP.Port)

	log.Info().Msgf("starting server on %s", serveAddr)

	if err := app.Listen(serveAddr); err != nil {
		return errors.Wrap(err, "fiber server run failed")
	}

	return nil
}

func main() {
	config.InitLogging()
	if err := run(); err != nil {
		fmt.Printf("%+v\n", err)
		log.Error().Stack().Err(err).Msg("failed to run app")
	}
}
