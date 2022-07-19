package main

import (
	"embed"
	"fmt"
	"github.com/codemicro/wiki/wiki/config"
	"github.com/codemicro/wiki/wiki/db"
	"github.com/codemicro/wiki/wiki/endpoints"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

//go:generate sass static/:static/

//go:embed static
var staticAssets embed.FS

func run() error {
	err := config.SAML.Load()
	if err != nil {
		return errors.WithStack(err)
	}

	database, err := db.New()
	if err != nil {
		return errors.WithStack(err)
	}

	if err := database.Migrate(); err != nil {
		return errors.Wrap(err, "failed migration")
	}

	e, err := endpoints.New(database)
	if err != nil {
		return errors.WithStack(err)
	}

	app := e.SetupApp()

	app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(staticAssets),
		PathPrefix: "static",
	}))

	serveAddr := config.HTTP.Host + ":" + strconv.Itoa(config.HTTP.Port)

	log.Info().Msgf("starting server on http://%s", serveAddr)

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
