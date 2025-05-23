package cmd

import (
	"github.com/SticketInya/kredentials/kredentials"
	"github.com/spf13/cobra"
)

func NewRootCmd(cli *kredentials.KredentialsCli) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "kredentials",
		Short: "kredentials is a cli tool for managing kubernetes configs",
		Long:  "kredentials is a cli tool for managing kubernetes configs",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

	rootCmd.AddCommand(NewAddCmd(cli))
	rootCmd.AddCommand(NewListCmd(cli))
	rootCmd.AddCommand(NewUseCommand(cli))
	rootCmd.AddCommand(NewDeleteCmd(cli))
	rootCmd.AddCommand(NewVersionCmd(cli))
	rootCmd.AddCommand(NewRevertCmd(cli))
	rootCmd.AddCommand(NewBackupCmd(cli))
	rootCmd.AddCommand(NewRestoreCmd(cli))

	return rootCmd
}
