package main

import (
	"fmt"
	"os"

	"github.com/SticketInya/kredentials/cmd"
	"github.com/SticketInya/kredentials/kredentials"
)

func main() {
	cli := kredentials.NewKredentialsCli(kredentials.NewKredentialsDefaultConfig())
	rootCmd := cmd.NewRootCmd(cli)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oopsie daisy! An error while executing kredentials '%s'\n", err)
		os.Exit(1)
	}
}
