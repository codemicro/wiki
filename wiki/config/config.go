package config

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

func InitLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

var HTTP = struct {
	Host string
	Port int
}{
	Host: asString(withDefault(fetchFromFile("http.host"), "0.0.0.0")),
	Port: asInt(withDefault(fetchFromFile("http.port"), 8080)),
}

var Database = struct {
	Filename string
}{
	Filename: asString(withDefault(fetchFromFile("database.filename"), "wiki.sqlite3.db")),
}
