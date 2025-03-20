package cmd

import (
	"fmt"

	"github.com/SticketInya/kredentials/internal/cmdutil"
	"github.com/SticketInya/kredentials/kredentials"
	"github.com/spf13/cobra"
)

func NewUseCommand(cli *kredentials.KredentialsCli) *cobra.Command {
	useCmd := &cobra.Command{
		Use:     "use name",
		Short:   "Use the selected kredential as kubernetes config",
		Args:    cobra.ExactArgs(1),
		PreRunE: validateUseArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUse(cli, args)
		},
		Example: "kredentials use my-kredential",
	}

	return useCmd
}

func runUse(cli *kredentials.KredentialsCli, args []string) error {
	kredentialName := args[0]
	if err := cli.Manager.UseKredential(kredentialName); err != nil {
		return fmt.Errorf("setting active kubernetes config: %w", err)
	}

	cli.Printer.Printf("now using '%s' as kubernetes config!\n", kredentialName)
	return nil
}

func validateUseArgs(cmd *cobra.Command, args []string) error {
	switch {
	case args[0] == "":
		return cmdutil.ErrWithUsage(cmd, "required 'name' cannot be empty")
	default:
		return nil
	}
}
