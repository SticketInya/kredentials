package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kredentials",
	Short: "kredentials is a cli tool for managing kubernetes configs",
	Long:  "kredentials is a cli tool for managing kubernetes configs",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error while executing kredentials '%s'\n", err)
		os.Exit(1)
	}
}
