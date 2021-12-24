package cmd

import (
	"flag"

	"k8s.io/klog/v2"

	"github.com/qqbuby/kuberos/pkg/cmd/cert"
	"github.com/qqbuby/kuberos/pkg/cmd/serve"

	"github.com/spf13/cobra"
)

func NewKuberosCommand() *cobra.Command {
	var cmds = &cobra.Command{
		Use: "kuberos",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	flags := cmds.Flags()
	loggingFlags := &flag.FlagSet{}
	klog.InitFlags(loggingFlags)
	flags.AddGoFlagSet(loggingFlags)

	cmds.AddCommand(serve.NewCmdServer())
	cmds.AddCommand(cert.NewCmdCert())

	return cmds
}
