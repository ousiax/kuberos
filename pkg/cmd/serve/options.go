package serve

import (
	"crypto/tls"
	"net/http"
)

type webhookOptions struct {
	certFile string
	keyFile  string
	port     int

	tlsConfig *tls.Config
	server    *http.Server
}

func (o *webhookOptions) configTLS() error {
	sCert, err := tls.LoadX509KeyPair(o.certFile, o.keyFile)
	if err != nil {
		return err
	}
	o.tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{sCert},
		// TODO: uses mutual tls after we agree on what cert the apiserver should use.
		// ClientAuth:   tls.RequireAndVerifyClientCert,
	}
	return nil
}
