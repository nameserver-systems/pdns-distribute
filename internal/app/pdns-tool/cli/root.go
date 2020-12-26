package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pdns-tool",
	Short: "A useful operation utility for pdns-distribute",
	Long:  `A useful operation utility for pdns-distribute. For example you can trigger a global resync of a zone, ...`,
	Args:  cobra.ExactArgs(0),
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
