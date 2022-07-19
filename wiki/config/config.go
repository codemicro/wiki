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
	SecureCookies     bool
}{
	Host:              asString(withDefault("http.host", "0.0.0.0")),
	Port:              asInt(withDefault("http.port", 8080)),
	TrustProxyHeaders: asBool(fetchFromFile("http.trustProxyHeaders")),
	ExternalURL:       asString(withDefault("http.externalURL", "https://localhost")),
	SecureCookies:     asBool(withDefault("http.secureCookies", true)),
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
	NameMappingKey  string
	EmailMappingKey string
}

var SAML = &samlConfig{
	Autoload:        asBool(fetchFromFile("saml.autoload")),
	MetadataURL:     asString(fetchFromFile("saml.metadataURL")),
	EntityID:        asString(fetchFromFile("saml.idp.entityID")),
	SSOURL:          asString(fetchFromFile("saml.idp.ssoURL")),
	rawIDPCert:      asString(fetchFromFile("saml.idp.signingCertificate")),
	NameMappingKey:  asString(withDefault("saml.mappings.name", "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/name")),
	EmailMappingKey: asString(withDefault("saml.mappings.email", "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress")),
}
