package cmd

import (
	"fmt"

	"github.com/SticketInya/kredentials/kredentials"
	"github.com/spf13/cobra"
)

func NewDeleteCmd(cli *kredentials.KredentialsCli) *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"remove", "rm"},
		Short:   "delete a kredential",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelete(cli, args)
		},
	}

	return deleteCmd
}

func runDelete(cli *kredentials.KredentialsCli, args []string) error {
	configName := args[0]
	if configName == "" {
		return fmt.Errorf("kredential name cannot be empty")
	}

	cli.Manager.DeleteKredential(configName)
	cli.Printer.Printf("kredential '%s' deleted!\n", configName)
	return nil
}
