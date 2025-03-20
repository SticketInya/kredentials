package cmd

import (
	"github.com/SticketInya/kredentials/internal/cmdutil"
	"github.com/SticketInya/kredentials/kredentials"
	"github.com/spf13/cobra"
)

func NewRestoreCmd(cli *kredentials.KredentialsCli) *cobra.Command {
	restoreCmd := &cobra.Command{
		Use:     "restore path",
		Short:   "Restore kredentials from a backup",
		Args:    cobra.ExactArgs(1),
		PreRunE: validateRestoreArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRestore(cli, args)
		},
		Example: "kredentials restore ./my-backup.zip",
	}

	return restoreCmd
}

func runRestore(cli *kredentials.KredentialsCli, args []string) error {
	path := args[0]
	if err := cli.Manager.RestoreKredentialBackup(path); err != nil {
		return err
	}

	cli.Printer.Printf("restored backup from '%s'!\n", path)
	return nil
}

func validateRestoreArgs(cmd *cobra.Command, args []string) error {
	switch {
	case args[0] == "":
		return cmdutil.ErrWithUsage(cmd, "required 'path' cannot be empty")
	default:
		return nil
	}
}
