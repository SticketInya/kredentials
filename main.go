package main

import (
	"os"

	"github.com/SticketInya/kredentials/cmd"
	"github.com/SticketInya/kredentials/internal/cmdutil"
	"github.com/SticketInya/kredentials/kredentials"
)

// Set at build time
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

func main() {
	cmdErrorHandler := cmdutil.NewCmdErrorHandler(os.Stderr)

	config, err := kredentials.NewKredentialsDefaultConfig(kredentials.NewVersionConfig(Version, Commit, BuildDate))
	cmdErrorHandler.HandleAndExit(err, 1)

	cli := kredentials.NewKredentialsCli(config)
	cmdErrorHandler.HandleAndExit(err, 1)

	rootCmd := cmd.NewRootCmd(cli)
	cmdErrorHandler.HandleAndExit(rootCmd.Execute(), 1)
}
