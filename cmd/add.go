package cmd

import (
	"fmt"

	"github.com/SticketInya/kredentials/kredentials"
	"github.com/SticketInya/kredentials/models"
	"github.com/spf13/cobra"
)

func NewAddCmd(cli *kredentials.KredentialsCli) *cobra.Command {
	options := models.AddKredentialOptions{}
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "add a new config file",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdd(cli, options, args)
		},
	}

	addCmd.Flags().BoolVar(&options.OverwriteExisting, "force", false, "Overwrite existing kredential entry on conflict")

	return addCmd
}

func runAdd(cli *kredentials.KredentialsCli, options models.AddKredentialOptions, args []string) error {
	name, path := args[0], args[1]

	if name == "" || path == "" {
		return fmt.Errorf("name or path is missing")
	}

	if err := cli.Manager.AddKredential(name, path, options); err != nil {
		return err
	}

	cli.Printer.Printf("kredential '%s' added!\n", name)
	return nil
}
