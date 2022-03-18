package cli

import (
	"github.com/spf13/cobra"
	"go.xargs.dev/bindl/command"
)

var bindlGenerateIgnorePath = ".gitignore"

var BindlGenerateIgnore = &cobra.Command{
	Use:   "ignore",
	Short: "Generate ignore file for bindl programs",
	Long: `Generate ignore file for bindl programs 

By default, Bindl will take ".gitignore" as input and append 
<output directory>/* if it doesn't already exist.

For example, with default output directory "bin":

  $ bindl ignore -f .gitignore
  $ tail -n 1 .gitignore
  bin/*

Supports typical ignore files. e.g. .dockerignore`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return command.UpdateIgnoreFile(defaultConfig, bindlGenerateIgnorePath)
	},
}

func init() {
	BindlGenerateIgnore.Flags().StringVarP(&bindlGenerateIgnorePath, "path", "p", bindlGenerateIgnorePath, "path to ignore file")
}
