package cmdutil

import (
	"fmt"

	"github.com/spf13/cobra"
)

func GetStandardUsageString(cmd *cobra.Command) string {
	return fmt.Sprintf("See '%s -h' for help and examples.", cmd.CommandPath())
}

func ErrWithUsage(cmd *cobra.Command, format string, a ...any) error {
	return fmt.Errorf("%w\n%s", fmt.Errorf(format, a...), GetStandardUsageString(cmd))
}
