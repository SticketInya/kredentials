package cmd

import (
	"sort"
	"strings"

	"github.com/SticketInya/kredentials/formatter/templates"
	"github.com/SticketInya/kredentials/kredentials"
	"github.com/SticketInya/kredentials/models"
	"github.com/spf13/cobra"
)

func NewListCmd(cli *kredentials.KredentialsCli) *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "list all the kredentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(cli, args)
		},
	}

	return listCmd
}

const defaultListSeparator string = ","

func runList(cli *kredentials.KredentialsCli, args []string) error {
	kreds, err := cli.Manager.ListKredentials()
	if err != nil {
		return err
	}

	if len(kreds) == 0 {
		cli.Printer.Println("No kredentials found")
		return nil
	}

	nodeList := buildKredentialNodeList(kreds)
	return cli.Printer.StructuredPrint(nodeList)
}

func buildKredentialNodeList(kreds []*models.Kredential) templates.KredentialNodeListTemplate {
	var data templates.KredentialNodeListTemplate

	for _, kred := range kreds {
		data.Items = append(data.Items, templates.KredentialNodeTemplate{
			Name:     kred.Name,
			Clusters: joinMapKeysGeneric(kred.Config.Clusters, defaultListSeparator),
			Contexts: joinMapKeysGeneric(kred.Config.Contexts, defaultListSeparator),
		})
	}

	return data
}

func joinMapKeysGeneric[T any](m map[string]T, separator string) string {
	if len(m) == 0 {
		return ""
	}

	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	return strings.Join(keys, separator)
}
