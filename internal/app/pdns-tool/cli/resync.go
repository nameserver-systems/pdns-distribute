//nolint:gochecknoinits
package cli

import (
	"github.com/spf13/cobra"
)

var resyncCmd = &cobra.Command{
	Use:   "resync",
	Short: "Triggers a resync of some resources",
	Long:  `Triggers a resync of some resources`,
	Args:  cobra.ExactArgs(0),
}

func init() {
	rootCmd.AddCommand(resyncCmd)
}
