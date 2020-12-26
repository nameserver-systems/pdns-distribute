//nolint:gochecknoinits
package cli

import (
	"github.com/spf13/cobra"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/app/pdns-tool/resync/zone"
)

var resynczonecmd = &cobra.Command{
	Use:   "zone [zonename]",
	Short: "Triggers a global resync of a zone",
	Long:  `Triggers a global resync of a zone`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		zone.Execute(args[0])
	},
}

func init() {
	resyncCmd.AddCommand(resynczonecmd)
}
