package serve

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"

	_ "k8s.io/api/admission/v1"
	"k8s.io/klog/v2"

	"github.com/qqbuby/kuberos/pkg/cmd/serve/admission/controller"
	cmduitl "github.com/qqbuby/kuberos/pkg/cmd/util"
)

const (
	flagCertFile = "tls-cert-file"
	flagKeyFile  = "tls-private-key-file"
	flagPort     = "port"
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
		"Secure port that the webhook listens on")

	return cmd
}

func (o *webhookOptions) Complete() error {
	err := o.configTLS()
	if err != nil {
		return err
	}
	o.server = &http.Server{
		Addr:      fmt.Sprintf(":%d", o.port),
		TLSConfig: o.tlsConfig,
	}

	http.HandleFunc("/image/policy", controller.ServeImagePolicy)
	http.HandleFunc("/livez", func(w http.ResponseWriter, req *http.Request) { w.Write([]byte("ok")) })
	http.HandleFunc("/readyz", func(w http.ResponseWriter, req *http.Request) { w.Write([]byte("ok")) })

	return nil
}

func (o *webhookOptions) Validate(cmd *cobra.Command) error {
	return nil
}

func (o *webhookOptions) Run() error {
	klog.V(2).Infof("Start serveing: %s", o.server.Addr)
	return o.server.ListenAndServeTLS("", "")
}
