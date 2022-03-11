package cli

import (
	"github.com/spf13/cobra"
	"go.xargs.dev/bindl/internal/log"
)

var All = []*cobra.Command{
	BindlInit,
	BindlGet,
}

var logLevel = "info"

var Root = &cobra.Command{
	Use:  "bindl",
	Long: "Bindl is a static binary downloader for project development and infrastructure.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if logLevel != "" {
			return log.SetLevel(logLevel)
		}
		return nil
	},
}

type CLIConfig struct {
	Path string
}

var defaultConfig = &CLIConfig{
	Path: "./bindl.yaml",
}

func init() {
	Root.PersistentFlags().StringVar(&logLevel, "log-level", logLevel, "Log level: trace, debug, info, disabled")
	Root.PersistentFlags().StringVar(&defaultConfig.Path, "config", defaultConfig.Path, "Path to configuration file.")
	Root.AddCommand(All...)
}
