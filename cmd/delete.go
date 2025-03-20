package cmd

import (
	"github.com/SticketInya/kredentials/internal/cmdutil"
	"github.com/SticketInya/kredentials/kredentials"
	"github.com/spf13/cobra"
)

func NewDeleteCmd(cli *kredentials.KredentialsCli) *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:     "delete name",
		Aliases: []string{"remove", "rm"},
		Short:   "Delete a kredential",
		Args:    cobra.ExactArgs(1),
		PreRunE: validateDeleteArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelete(cli, args)
		},
		Example: "kredential delete my-kredential",
	}

	return deleteCmd
}

func runDelete(cli *kredentials.KredentialsCli, args []string) error {
	configName := args[0]
	cli.Manager.DeleteKredential(configName)
	cli.Printer.Printf("kredential '%s' deleted!\n", configName)
	return nil
}

func validateDeleteArgs(cmd *cobra.Command, args []string) error {
	switch {
	case args[0] == "":
		return cmdutil.ErrWithUsage(cmd, "required 'name' cannot be empty")
	default:
		return nil
	}
}
