package cert

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
	flagCommonName    = "cn"
	flagGroups        = "group"
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
}

func NewCmdCert() *cobra.Command {
	o := CertOptions{}

	cmd := &cobra.Command{
		Use:   "cert",
		Short: "Create self-signed certificate.",
		Run: func(cmd *cobra.Command, args []string) {
			cmduitl.CheckErr(o.Complete())
			cmduitl.CheckErr(o.Validate())
			cmduitl.CheckErr(o.Run())
		},
	}

	cmd.Flags().StringVar(&o.commonName, flagCommonName, "", "CommonName")
	cmd.MarkFlagRequired(flagCommonName)
	cmd.Flags().StringArrayVar(&o.organization, flagGroups, nil, "Organization")
	cmd.MarkFlagRequired(flagGroups)

	cmd.Flags().StringVar(&keyOut, flagKeyOut, "",
		"This specifies the output filename to write a key to or standard output if this option is not specified.")

	cmd.Flags().StringVar(&certOut, flagCertOut, "",
		"Outputs the certificate block in PEM format. - default stdout")

	return cmd
}

func (o *CertOptions) Validate() error {
	return nil
}

func (o *CertOptions) Complete() error {
	return nil
}

func (o *CertOptions) Run() error {
	key, cert, err := cmduitlpkix.CreateSelfSignedCertificate(o.commonName, o.organization, nil)

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
