package cmd

import (
	"fmt"
	"os"

	"github.com/SticketInya/kredentials/formatter"
	"github.com/SticketInya/kredentials/kredentials"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"remove", "rm"},
	Short:   "delete a kredential",
	Args:    cobra.ExactArgs(1),
	RunE:    runDelete,
}

func runDelete(cmd *cobra.Command, args []string) error {
	configName := args[0]
	if configName == "" {
		return fmt.Errorf("kredential name cannot be empty")
	}

	if !kredentials.CheckKredentialInStorage(configName) {
		return fmt.Errorf("kredential '%s' does not exists", configName)
	}

	if err := kredentials.DeleteKredentialFromStorage(configName); err != nil {
		return fmt.Errorf("deleting kredential '%s' %w", configName, err)
	}

	w := formatter.NewStructuredPrinter(os.Stdout)
	w.Printf("kredential '%s' deleted!\n", configName)
	return nil
}
