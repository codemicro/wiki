package config

import (
	"bytes"
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/ucarion/saml"
	"io/ioutil"
	"net/http"
	"strings"
)

const samlHTTPPostSSOBinding = "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"

func (sc *samlConfig) Load() error {
	if sc.Autoload {
		log.Info().Msg("Automatically loading SAML SSO metadata")
		return sc.autoloadSAMLMetadata()
	}

	cert, err := loadBase64Certificate([]byte(sc.rawIDPCert))
	if err != nil {
		return errors.WithStack(err)
	}
	sc.IDPCert = cert

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

	bodyContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.WithStack(err)
	}

	_ = ioutil.WriteFile("metadata", bodyContent, 0777)

	metadata := new(saml.EntityDescriptor)
	if err := xml.Unmarshal(bodyContent, metadata); err != nil {
		return errors.WithStack(err)
	}

	sc.EntityID = metadata.EntityID

	for _, ssoService := range metadata.IDPSSODescriptor.SingleSignOnServices {
		if strings.EqualFold(ssoService.Binding, samlHTTPPostSSOBinding) {
			sc.SSOURL = ssoService.Location
			goto foundBinding
		}
	}
	return errors.WithStack(fmt.Errorf("no compatible SSO service for binding %s", samlHTTPPostSSOBinding))

foundBinding:
	sc.IDPCert, err = loadBase64Certificate(
		[]byte(metadata.IDPSSODescriptor.KeyDescriptor.KeyInfo.X509Data.X509Certificate.Value),
	)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func loadBase64Certificate(encoded []byte) (*x509.Certificate, error) {
	decoded := make([]byte, base64.RawStdEncoding.DecodedLen(len(encoded)))
	_, err := base64.StdEncoding.Decode(decoded, encoded)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	decoded = bytes.TrimRight(decoded, "\x00")

	cert, err := x509.ParseCertificate(decoded)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return cert, nil
}
