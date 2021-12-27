package serve

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"

	"k8s.io/klog/v2"

	cmduitl "github.com/qqbuby/kuberos/pkg/cmd/util"
	"github.com/qqbuby/kuberos/pkg/plugins/pods"
)

func NewCmdServer() *cobra.Command {
	o := &webhookOptions{}
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Kubernetes Extension Admission Control.",
		Long:  "Kubernetes Extension Admission Control.",
		Run: func(cmd *cobra.Command, args []string) {
			cmduitl.CheckErr(o.Complete())
			cmduitl.CheckErr(o.Validate(cmd))
			cmduitl.CheckErr(o.Run())
		},
	}

	cmd.Flags().StringVar(&o.certFile, flagCertFile, "",
		"File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated after server cert).")
	cmd.MarkFlagRequired(flagCertFile)

	cmd.Flags().StringVar(&o.keyFile, flagKeyFile, "",
		"File containing the default x509 private key matching --tls-cert-file.")
	cmd.MarkFlagRequired(flagKeyFile)

	cmd.Flags().IntVar(&o.port, flagPort, 443,
		"Secure port that the webhook listens on.")

	return cmd
}

const (
	flagCertFile = "tls-cert-file"
	flagKeyFile  = "tls-private-key-file"
	flagPort     = "port"
)

type webhookOptions struct {
	certFile string
	keyFile  string
	port     int

	mux       *http.ServeMux
	tlsConfig *tls.Config
	server    *http.Server
}

func (o *webhookOptions) configServer() error {
	sCert, err := tls.LoadX509KeyPair(o.certFile, o.keyFile)
	if err != nil {
		return err
	}
	o.tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{sCert},
		// TODO: uses mutual tls after we agree on what cert the apiserver should use.
		// ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	o.mux = http.NewServeMux()
	o.server = &http.Server{
		Handler:   o.mux,
		Addr:      fmt.Sprintf(":%d", o.port),
		TLSConfig: o.tlsConfig,
	}

	return nil
}

func (o *webhookOptions) handleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	if handler == nil {
		panic("http: nil handler")
	}
	o.mux.HandleFunc(pattern, handler)
}

func (o *webhookOptions) Complete() error {
	err := o.configServer()
	if err != nil {
		return err
	}

	o.handleFunc("/livez", func(w http.ResponseWriter, req *http.Request) { w.Write([]byte("ok")) })
	o.handleFunc("/readyz", func(w http.ResponseWriter, req *http.Request) { w.Write([]byte("ok")) })

	o.handleFunc("/api/v1/pods", pods.ServeHTTP)

	return nil
}

func (o *webhookOptions) Validate(cmd *cobra.Command) error {
	return nil
}

func (o *webhookOptions) Run() error {
	klog.V(2).Infof("Start serveing: %s", o.server.Addr)
	return o.server.ListenAndServeTLS("", "")
}
