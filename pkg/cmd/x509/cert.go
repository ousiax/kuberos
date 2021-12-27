package x509

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	cmduitl "github.com/qqbuby/kuberos/pkg/cmd/util"
	cmduitlpkix "github.com/qqbuby/kuberos/pkg/cmd/util/pkix"
)

const (
	flagKeyOut        = "keyout"
	flagCertOut       = "out"
	flagCommonName    = "common-name"
	flagOrganization  = "organization"
	flagDnsNames      = "dns-names"
	flagExpiration    = "expiration"
	expirationSeconds = 60 * 60 * 24 * 365 // one year in seconds
)

var (
	keyOut  string
	certOut string
)

type CertOptions struct {
	commonName   string
	organization []string
	dnsNames     []string
}

func NewCmdCert() *cobra.Command {
	o := CertOptions{}

	cmd := &cobra.Command{
		Use:   "self-signed-cert",
		Short: "Create self-signed certificate.",
		Run: func(cmd *cobra.Command, args []string) {
			cmduitl.CheckErr(o.Complete(cmd, args))
			cmduitl.CheckErr(o.Validate())
			cmduitl.CheckErr(o.Run())
		},
	}

	cmd.Flags().StringVarP(&o.commonName, flagCommonName, "u", "", "CommonName")
	cmd.Flags().StringArrayVarP(&o.organization, flagOrganization, "g", nil, "Organization")
	cmd.Flags().StringArrayVarP(&o.dnsNames, flagDnsNames, "n", nil, "Subject Alternate DNS Names.")

	cmd.Flags().StringVarP(&keyOut, flagKeyOut, "k", "",
		"This specifies the output filename to write the private key in PKCS #8 PEM format to or standard output if this option is not specified.")

	cmd.Flags().StringVarP(&certOut, flagCertOut, "o", "",
		"This specifies the output filename to write the certificate in PEM format to or standard output if this option is not specified.")

	return cmd
}

func (o *CertOptions) Complete(cmd *cobra.Command, args []string) error {
	if cn := cmduitl.GetFlagString(cmd, flagCommonName); len(cn) == 0 {
		cmd.Usage()
		os.Exit(1)

	}
	return nil
}

func (o *CertOptions) Validate() error {
	return nil
}

func (o *CertOptions) Run() error {
	key, cert, err := cmduitlpkix.CreateSelfSignedCertificate(o.commonName, o.organization, o.dnsNames)

	if err != nil {
		return err
	}

	keyPem, err := cmduitlpkix.PemPkcs8PKey(key)
	if err != nil {
		return err
	}
	if len(keyOut) > 0 {
		if err := os.WriteFile(keyOut, keyPem, 0400); err != nil {
			return err
		}
	} else {
		fmt.Fprint(os.Stdout, string(keyPem))
	}

	certPem, err := cmduitlpkix.PemCertificate(cert)
	if err != nil {
		return err
	}
	if len(certOut) > 0 {
		if err := os.WriteFile(certOut, certPem, 0400); err != nil {
			return err
		}
	} else {
		fmt.Fprint(os.Stdout, string(certPem))
	}

	return nil
}
