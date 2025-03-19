package cmd

import (
	"fmt"

	"github.com/SticketInya/kredentials/kredentials"
	"github.com/spf13/cobra"
)

func NewRevertCmd(cli *kredentials.KredentialsCli) *cobra.Command {
	revertCmd := &cobra.Command{
		Use:   "revert",
		Short: "Revert the configuration in use to the previous one",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRevert(cli)
		},
	}

	return revertCmd
}

func runRevert(cli *kredentials.KredentialsCli) error {
	if err := cli.Manager.RevertKredential(); err != nil {
		return fmt.Errorf("cannot revert kredential: %w", err)
	}

	cli.Printer.Println("kredential reverted!")
	return nil
}
