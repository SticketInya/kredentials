package cmd

import (
	"github.com/SticketInya/kredentials/internal/cmdutil"
	"github.com/SticketInya/kredentials/kredentials"
	"github.com/spf13/cobra"
)

func NewBackupCmd(cli *kredentials.KredentialsCli) *cobra.Command {
	backupCmd := &cobra.Command{
		Use:     "backup path",
		Short:   "Create a backup of your kredentials",
		Example: "kredentials backup ./my-backup-dir",
		Args:    cobra.ExactArgs(1),
		PreRunE: validateBackupArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBackup(cli, args)
		},
	}

	return backupCmd
}

func runBackup(cli *kredentials.KredentialsCli, args []string) error {
	path := args[0]
	if err := cli.Manager.CreateKredentialBackup(path); err != nil {
		return err
	}

	cli.Printer.Printf("created backup at '%s'!\n", path)
	return nil
}

func validateBackupArgs(cmd *cobra.Command, args []string) error {
	switch {
	case args[0] == "":
		return cmdutil.ErrWithUsage(cmd, "required 'path' cannot be empty")
	default:
		return nil
	}
}
