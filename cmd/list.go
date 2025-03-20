package cmd

import (
	"github.com/SticketInya/kredentials/formatter/templates"
	"github.com/SticketInya/kredentials/kredentials"
	"github.com/spf13/cobra"
)

func NewListCmd(cli *kredentials.KredentialsCli) *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all the kredentials",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(cli)
		},
	}

	return listCmd
}

func runList(cli *kredentials.KredentialsCli) error {
	kreds, err := cli.Manager.ListKredentials()
	if err != nil {
		return err
	}

	if len(kreds) == 0 {
		cli.Printer.Println("no kredentials found")
		return nil
	}

	nodeList := templates.BuildKredentialNodeList(kreds)
	return cli.Printer.StructuredPrint(nodeList)
}
