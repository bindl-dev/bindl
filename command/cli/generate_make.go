package cli

import (
	"github.com/spf13/cobra"
	"go.xargs.dev/bindl/command"
)

var bindlGenerateMakefilePath = "Makefile.bindl"

var BindlGenerateMake = &cobra.Command{
	Use:   "make",
	Short: "Generate Makefile for bindl programs",
	Long: `Generate Makefile for all programs in lockfile.

By default, the generated Makefile will be named 'Makefile.bindl', which can be
imported by the project's primary Makefile using 'include' directive.

After including, you can use rules defined in 'Makefile.bindl' as a dependency
in your other rules. For example:

	$ head -n 5 Makefile
	include Makefile.bindl

	.PHONY: container
	container: bin/ko
		bin/ko publish -B .

Calling the imported rules also works on 'make' CLI.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return command.GenerateMakefile(defaultConfig, bindlGenerateMakefilePath)
	},
}

func init() {
	BindlGenerateMake.Flags().StringVarP(&bindlGenerateMakefilePath, "path", "p", bindlGenerateMakefilePath, "path to generated Makefile")
}
