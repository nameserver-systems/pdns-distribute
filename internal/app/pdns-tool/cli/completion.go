//nolint:gochecknoinits
package cli

import (
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
)

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates bash completion scripts",
	Long: `To load completion run

. <(bitbucket completion)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(bitbucket completion)
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := rootCmd.GenBashCompletion(os.Stdout)
		if err != nil {
			logger.FatalErrLog(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
