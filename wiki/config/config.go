package config

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	dsig "github.com/russellhaering/goxmldsig"
)

func InitLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

var HTTP = struct {
	Host              string
	Port              int
	TrustProxyHeaders bool
	ExternalURL       string
}{
	Host:              asString(withDefault("http.host", "0.0.0.0")),
	Port:              asInt(withDefault("http.port", 8080)),
	TrustProxyHeaders: asBool(fetchFromFile("http.trustProxyHeaders")),
	ExternalURL:       asString(withDefault("http.externalURL", "https://localhost")),
}

var Database = struct {
	Filename string
}{
	Filename: asString(withDefault("database.filename", "wiki.sqlite3.db")),
}

type samlConfig struct {
	Autoload        bool
	MetadataURL     string
	EntityID        string
	SSOURL          string
	rawIDPCert      string
	IDPCertificates dsig.X509CertificateStore
}

var SAML = &samlConfig{
	Autoload:    asBool(fetchFromFile("saml.autoload")),
	MetadataURL: asString(fetchFromFile("saml.metadataURL")),
	EntityID:    asString(fetchFromFile("saml.idp.entityID")),
	SSOURL:      asString(fetchFromFile("saml.idp.ssoURL")),
	rawIDPCert:  asString(fetchFromFile("saml.idp.signingCertificate")),
}
