package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.xargs.dev/bindl/config"
)

var bindlListOneline = false

var BindlList = &cobra.Command{
	Use:   "list",
	Short: "List out programs defined in lockfile",
	Long:  "List shows all program names defined in lockfile.",

	Args: cobra.NoArgs,

	RunE: func(cmd *cobra.Command, args []string) error {
		l, err := config.ParseLock(defaultConfig.LockfilePath)
		if err != nil {
			return fmt.Errorf("parsing lockfile: %w", err)
		}
		separator := "\n"
		if bindlListOneline {
			separator = " "
		}
		for _, p := range l.Programs {
			fmt.Printf("%s%s", p.PName, separator)
		}
		if bindlListOneline {
			fmt.Println()
		}
		return nil
	},
}

func init() {
	BindlList.Flags().BoolVar(&bindlListOneline, "oneline", bindlListOneline, "List programs in one line, space separated")
}
