package cli

import (
	"github.com/spf13/cobra"
	"go.xargs.dev/bindl/command"
)

var bindlIgnorePath = ".gitignore"

var BindlIgnore = &cobra.Command{
	Use:   "ignore",
	Short: "Generate ignore file for bindl programs",
	Long: `Generate ignore file for bindl programs 

By default, Bindl will take ".gitignore" as input and append 
<output directory>/*.

For example, with default output directory "bin":

  $ bindl ignore -f .gitignore
  $ tail -n 1 .gitignore
  bin/*

This operation is idempotent, so running the command multiple times will
result in the same ignore file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return command.UpdateIgnoreFile(defaultConfig, bindlIgnorePath)
	},
}

func init() {
	BindlIgnore.Flags().StringVarP(&bindlIgnorePath, "path", "p", bindlIgnorePath, "path to ignore file")
}
