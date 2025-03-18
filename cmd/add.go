package cmd

import (
	"fmt"
	"os"

	"github.com/SticketInya/kredentials/kredentials"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a new config file",
	Args:  cobra.ExactArgs(2),
	RunE:  runAdd,
}

func runAdd(cmd *cobra.Command, args []string) error {
	name, path := args[0], args[1]

	if name == "" || path == "" {
		return fmt.Errorf("name or path is missing")
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("file '%s' does not exists", path)
	}

	// TODO: add dir support
	if info.IsDir() {
		return fmt.Errorf("invalid path provided, '%s' is a directory", path)
	}

	kred := kredentials.NewKredential(name)
	if err = kred.ReadKubernetesConfig(path); err != nil {
		return fmt.Errorf("reading kubernetes config %w", err)
	}

	if err = kred.StoreKredential(kredentials.DefaultConfigStorageDir); err != nil {
		return err
	}
	fmt.Printf("'%s' config added!", name)

	return nil
}
