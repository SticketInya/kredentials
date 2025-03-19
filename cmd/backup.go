package cmd

import (
	"fmt"

	"github.com/SticketInya/kredentials/kredentials"
	"github.com/spf13/cobra"
)

func NewBackupCmd(cli *kredentials.KredentialsCli) *cobra.Command {
	backupCmd := &cobra.Command{
		Use:   "backup",
		Short: "Create a backup of your kredentials",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBackup(cli, args)
		},
	}

	return backupCmd
}

func runBackup(cli *kredentials.KredentialsCli, args []string) error {
	path := args[0]

	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	if err := cli.Manager.CreateKredentialBackup(path); err != nil {
		return err
	}

	cli.Printer.Printf("created backup at '%s'!\n", path)
	return nil
}
