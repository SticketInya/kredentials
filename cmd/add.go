package cmd

import (
	"github.com/SticketInya/kredentials/internal/cmdutil"
	"github.com/SticketInya/kredentials/kredentials"
	"github.com/SticketInya/kredentials/models"
	"github.com/spf13/cobra"
)

func NewAddCmd(cli *kredentials.KredentialsCli) *cobra.Command {
	options := models.AddKredentialOptions{}
	addCmd := &cobra.Command{
		Use:     "add name path",
		Short:   "Add a new kredential",
		Args:    cobra.ExactArgs(2),
		PreRunE: validateAddArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdd(cli, options, args)
		},
		Example: "kredentials add my-config ./kubeconfig",
	}

	addCmd.Flags().BoolVar(&options.OverwriteExisting, "force", false, "Overwrite existing kredential entry on conflict")

	return addCmd
}

func runAdd(cli *kredentials.KredentialsCli, options models.AddKredentialOptions, args []string) error {
	name, path := args[0], args[1]
	if err := cli.Manager.AddKredential(name, path, options); err != nil {
		return err
	}

	cli.Printer.Printf("kredential '%s' added!\n", name)
	return nil
}

func validateAddArgs(cmd *cobra.Command, args []string) error {
	switch {
	case args[0] == "":
		return cmdutil.ErrWithUsage(cmd, "required 'name' cannot be empty")
	case args[1] == "":
		return cmdutil.ErrWithUsage(cmd, "required 'path' cannot be empty")
	default:
		return nil
	}
}
