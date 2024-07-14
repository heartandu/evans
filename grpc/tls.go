package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/pkg/errors"
)

var ErrMutualAuthParamsAreNotEnough = errors.New("cert and certkey are required to authenticate mutually")

func tlsConfig(cacert, cert, certKey string) (*tls.Config, error) {
	var tlsCfg tls.Config

	if cacert != "" {
		b, err := os.ReadFile(cacert)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read the CA certificate")
		}

		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(b) {
			return nil, errors.New("failed to append the client certificate")
		}

		tlsCfg.RootCAs = cp
	}

	if cert != "" && certKey != "" {
		// Enable mutual authentication
		certificate, err := tls.LoadX509KeyPair(cert, certKey)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read the client certificate")
		}

		tlsCfg.Certificates = append(tlsCfg.Certificates, certificate)
	} else if cert != "" || certKey != "" {
		return nil, ErrMutualAuthParamsAreNotEnough
	}

	return &tlsCfg, nil
}
