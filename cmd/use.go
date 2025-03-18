package cmd

import (
	"fmt"
	"os"

	"github.com/SticketInya/kredentials/formatter"
	"github.com/SticketInya/kredentials/kredentials"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(useCmd)
}

var useCmd = &cobra.Command{
	Use:   "use",
	Short: "use the selected kubernetes config",
	Args:  cobra.ExactArgs(1),
	RunE:  runUse,
}

func runUse(cmd *cobra.Command, args []string) error {
	configName := args[0]

	if configName == "" {
		return fmt.Errorf("config name cannot be empty")
	}

	kred, err := kredentials.RetrieveKredentialFromStorage(configName)
	if err != nil {
		return fmt.Errorf("retrieving '%s' from storage %w", configName, err)
	}

	if err = kred.SetAsKubernetesConfig(); err != nil {
		return fmt.Errorf("setting '%s' as kubernetes config %w", configName, err)
	}

	w := formatter.NewStructuredPrinter(os.Stdout)
	w.Printf("Now using '%s' as kubernetes config!\n", kred.Name)
	return nil
}
