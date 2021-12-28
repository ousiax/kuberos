package cmd

import (
	"flag"

	"k8s.io/klog/v2"

	"github.com/qqbuby/kuberos/pkg/cmd/serve"
	"github.com/qqbuby/kuberos/pkg/cmd/x509"

	"github.com/spf13/cobra"
)

func NewKuberosCommand() *cobra.Command {
	var cmds = &cobra.Command{
		Use:     "kuberos",
		Version: "v0.1",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	flags := cmds.PersistentFlags()
	loggingFlags := &flag.FlagSet{}
	klog.InitFlags(loggingFlags)
	flags.AddGoFlagSet(loggingFlags)

	cmds.AddCommand(serve.NewCmdServer())
	cmds.AddCommand(x509.NewCmdCert())

	return cmds
}
