package cmd

import (
	"fmt"

	"github.com/SticketInya/kredentials/kredentials"
	"github.com/spf13/cobra"
)

func NewAddCmd(cli *kredentials.KredentialsCli) *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "add a new config file",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdd(cli, args)
		},
	}

	return addCmd
}

func runAdd(cli *kredentials.KredentialsCli, args []string) error {
	name, path := args[0], args[1]

	if name == "" || path == "" {
		return fmt.Errorf("name or path is missing")
	}

	if err := cli.Manager.AddKredential(name, path); err != nil {
		return err
	}

	cli.Printer.Printf("kredential '%s' added!\n", name)
	return nil
}
