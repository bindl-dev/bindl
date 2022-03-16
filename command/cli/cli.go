package cli

import (
	"github.com/spf13/cobra"
	"go.xargs.dev/bindl/config"
	"go.xargs.dev/bindl/internal/log"
)

var All = []*cobra.Command{
	BindlGet,
	BindlSync,
	BindlMake,
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
	Root.PersistentFlags().BoolVarP(&logDisable, "silent", "s", logDisable, "silence logs")
	Root.PersistentFlags().BoolVar(&logDebug, "debug", logDebug, "show debug logs")
	Root.PersistentFlags().StringVarP(&defaultConfig.Path, "config", "c", defaultConfig.Path, "path to configuration file")
	Root.PersistentFlags().StringVarP(&defaultConfig.LockfilePath, "lock", "l", defaultConfig.LockfilePath, "path to lockfile")
	Root.PersistentFlags().StringVarP(&defaultConfig.OutputDir, "bin", "b", defaultConfig.OutputDir, "directory to store binaries")
	Root.AddCommand(All...)
}
