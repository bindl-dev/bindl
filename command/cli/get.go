package cli

import (
	"github.com/spf13/cobra"

	"go.xargs.dev/bindl/command"
)

var bindlGetAll bool

var BindlGet = &cobra.Command{
	Use:   "get NAME",
	Short: "Get local copy of program",
	Long: `Get downloads the names program, which must already exist in bindl.yaml,
and ensures the program is ready to be used by setting executable flag.`,
	RunE: func(cmd *cobra.Command, names []string) error {
		if bindlGetAll {
			return command.GetAll(cmd.Context())
		} else {
			return command.Get(cmd.Context(), names...)
		}
	},
}

func init() {
	BindlGet.Flags().BoolVarP(&bindlGetAll, "all", "a", false, "Get all programs defined in 'bindl-lock.yaml'")
}
