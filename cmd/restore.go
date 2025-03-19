package cmd

import (
	"fmt"

	"github.com/SticketInya/kredentials/kredentials"
	"github.com/spf13/cobra"
)

func NewRestoreCmd(cli *kredentials.KredentialsCli) *cobra.Command {
	restoreCmd := &cobra.Command{
		Use:   "restore",
		Short: "restore kredentials from a backup",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRestore(cli, args)
		},
	}

	return restoreCmd
}

func runRestore(cli *kredentials.KredentialsCli, args []string) error {
	path := args[0]

	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	if err := cli.Manager.RestoreKredentialBackup(path); err != nil {
		return err
	}

	cli.Printer.Printf("restored backup from '%s'!\n", path)
	return nil
}
