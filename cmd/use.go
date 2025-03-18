package cmd

import (
	"fmt"

	"github.com/SticketInya/kredentials/kredentials"
	"github.com/spf13/cobra"
)

func NewUseCommand(cli *kredentials.KredentialsCli) *cobra.Command {
	useCmd := &cobra.Command{
		Use:   "use",
		Short: "use the selected kredential as kubernetes config",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUse(cli, args)
		},
	}

	return useCmd
}

func runUse(cli *kredentials.KredentialsCli, args []string) error {
	kredentialName := args[0]
	if kredentialName == "" {
		return fmt.Errorf("kredential name cannot be empty")
	}

	if err := cli.Manager.UseKredential(kredentialName); err != nil {
		return fmt.Errorf("setting active kubernetes config: %w", err)
	}

	cli.Printer.Printf("Now using '%s' as kubernetes config!\n", kredentialName)
	return nil
}
