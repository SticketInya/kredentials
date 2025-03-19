package main

import (
	"fmt"
	"os"

	"github.com/SticketInya/kredentials/cmd"
	"github.com/SticketInya/kredentials/kredentials"
)

// Set at build time
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

func main() {
	config, err := kredentials.NewKredentialsDefaultConfig(kredentials.NewVersionConfig(Version, Commit, BuildDate))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	cli := kredentials.NewKredentialsCli(config)
	if err = cli.Initialize(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	rootCmd := cmd.NewRootCmd(cli)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oopsie daisy! An error while executing kredentials '%s'\n", err)
		os.Exit(1)
	}
}
