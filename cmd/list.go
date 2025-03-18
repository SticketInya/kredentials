package cmd

import (
	"os"
	"sort"
	"strings"

	"github.com/SticketInya/kredentials/formatter"
	"github.com/SticketInya/kredentials/formatter/templates"
	"github.com/SticketInya/kredentials/kredentials"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

const defaultListSeparator string = ","

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all the kredentials",
	RunE:  runList,
}

func runList(cmd *cobra.Command, args []string) error {
	kreds, err := kredentials.ReadKredentials(kredentials.DefaultConfigStorageDir)
	if err != nil {
		// TODO: check if returning error here is more appropriate
		return nil
	}

	printer := formatter.NewStructuredPrinter(os.Stdout)

	if len(kreds) == 0 {
		printer.Println("No kredentials found")
		return nil
	}

	nodeList := buildKredentialNodeList(kreds)
	return printer.StructuredPrint(nodeList)
}

func buildKredentialNodeList(kreds []*kredentials.Kredential) templates.KredentialNodeListTemplate {
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
