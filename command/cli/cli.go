package cli

import (
	"github.com/spf13/cobra"
	"go.xargs.dev/bindl/config"
	"go.xargs.dev/bindl/internal/log"
)

var All = []*cobra.Command{
	BindlInit,
	BindlGet,
	BindlSync,
}

var logDebug bool
var logDisable bool

var Root = &cobra.Command{
	Use:  "bindl",
	Long: "Bindl is a static binary downloader for project development and infrastructure.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if logDisable {
			return log.SetLevel("disabled")
		} else if logDebug {
			return log.SetLevel("debug")
		} else {
			return log.SetLevel("info")
		}
	},
}

var defaultConfig = &config.Runtime{
	Path:         "./bindl.yaml",
	LockfilePath: "./.bindl-lock.yaml",
	OutputDir:    "./bin",
}

func init() {
	Root.PersistentFlags().BoolVar(&logDisable, "silent", logDisable, "Silence logs")
	Root.PersistentFlags().BoolVar(&logDebug, "debug", logDebug, "Show debug logs")
	Root.PersistentFlags().StringVar(&defaultConfig.Path, "config", defaultConfig.Path, "Path to configuration file.")
	Root.PersistentFlags().StringVar(&defaultConfig.LockfilePath, "lock", defaultConfig.LockfilePath, "Path to lockfile.")
	Root.PersistentFlags().StringVarP(&defaultConfig.OutputDir, "out", "o", defaultConfig.OutputDir, "Directory to store binaries.")
	Root.AddCommand(All...)
}
