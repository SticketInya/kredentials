package cmd

import (
	"github.com/SticketInya/kredentials/kredentials"
	"github.com/spf13/cobra"
)

func NewVersionCmd(cli *kredentials.KredentialsCli) *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the application version information",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			runVersion(cli)
		},
	}

	return versionCmd
}

func runVersion(cli *kredentials.KredentialsCli) {
	versionInformation := cli.GetVersion()
	cli.Printer.Printf(
		"Application Version: %s\nCommit Hash: %s\nBuild Date: %s\n",
		versionInformation.ApplicationVersion,
		versionInformation.CommitHash,
		versionInformation.BuildDate,
	)
}
