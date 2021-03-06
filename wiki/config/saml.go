package config

import (
	"bytes"
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	samlTypes "github.com/russellhaering/gosaml2/types"
	dsig "github.com/russellhaering/goxmldsig"
	"net/http"
	"strings"
)

const (
	httpPostBinding     = "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"
	httpRedirectBinding = "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect"
)

func (sc *samlConfig) Load() error {
	if sc.Autoload {
		if sc.MetadataURL == "" {
			log.Fatal().Msg("required key saml.metadataURL not found (required by saml.autoload=true)")
		}
		log.Info().Msg("automatically loading SAML SSO metadata")
		return sc.autoloadSAMLMetadata()
	}

	if sc.EntityID == "" {
		log.Fatal().Msg("required key saml.idp.entityID not found (required by saml.autoload=false)")
	} else if sc.SSOURL == "" {
		log.Fatal().Msg("required key saml.idp.ssoURL not found (required by saml.autoload=false)")
	} else if sc.rawIDPCert == "" {
		log.Fatal().Msg("required key saml.idp.signingCertificate not found (required by saml.autoload=false)")
	}

	cert, err := loadBase64Certificate([]byte(sc.rawIDPCert))
	if err != nil {
		return errors.WithStack(err)
	}
	sc.IDPCertificates = &dsig.MemoryX509CertificateStore{Roots: []*x509.Certificate{cert}}

	return nil
}

func (sc *samlConfig) autoloadSAMLMetadata() error {
	req, err := http.NewRequest(http.MethodGet, sc.MetadataURL, nil)
	if err != nil {
		return errors.WithStack(err)
	}

	req.Header.Add("User-Agent", "codemicro/wiki (+https://github.com/codemicro/wiki)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.WithStack(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("got non-200 status code for SAML metadata endpoint during autoload")
	}

	metadata := new(samlTypes.EntityDescriptor)
	if err := xml.NewDecoder(resp.Body).Decode(metadata); err != nil {
		return errors.WithStack(err)
	}

	sc.EntityID = metadata.EntityID

	for _, ssoService := range metadata.IDPSSODescriptor.SingleSignOnServices {
		if strings.EqualFold(ssoService.Binding, httpRedirectBinding) {
			sc.SSOURL = ssoService.Location
			goto foundBinding
		}
	}
	return errors.WithStack(fmt.Errorf("no compatible SSO service for binding %s", httpRedirectBinding))

foundBinding:

	certStore := &dsig.MemoryX509CertificateStore{
		Roots: []*x509.Certificate{},
	}

	for _, kd := range metadata.IDPSSODescriptor.KeyDescriptors {
		for _, xcert := range kd.KeyInfo.X509Data.X509Certificates {
			cert, err := loadBase64Certificate([]byte(xcert.Data))
			if err != nil {
				return errors.WithStack(err)
			}

			certStore.Roots = append(certStore.Roots, cert)
		}
	}

	sc.IDPCertificates = certStore

	return nil
}

func loadBase64Certificate(encoded []byte) (*x509.Certificate, error) {
	decoded := make([]byte, base64.RawStdEncoding.DecodedLen(len(encoded)))
	_, err := base64.StdEncoding.Decode(decoded, encoded)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Sometimes we get an extra null byte from YAML certificates
	decoded = bytes.TrimRight(decoded, "\x00")

	cert, err := x509.ParseCertificate(decoded)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return cert, nil
}
