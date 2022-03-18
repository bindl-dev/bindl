package cli

import "github.com/spf13/cobra"

var BindlGenerate = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "Generate project integration files (e.g. Makefile, .gitignore)",
	Long:    `Generate common setup for projects to have smooth workflow.`,
}

func init() {
	BindlGenerate.AddCommand(BindlGenerateIgnore)
	BindlGenerate.AddCommand(BindlGenerateMake)
}
