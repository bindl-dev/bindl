package cli

import (
	"github.com/spf13/cobra"
	"go.xargs.dev/bindl/command"
)

var bindlSyncStdout bool

var BindlSync = &cobra.Command{
	Use:   "sync",
	Short: "Sync configuration and lockfile",
	Long: `Synchronize bindl configuration with lockfile.

Sync will update lockfile (i.e. .bindl-lock.yaml) according to configuration
file specifications (i.e. bindl.yaml), ensuring that checksums exists in
lockfile for all desired platforms and programs.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return command.Sync(cmd.Context(), defaultConfig, bindlSyncStdout)
	},
}

func init() {
	BindlSync.Flags().BoolVar(&bindlSyncStdout, "stdout", false, "write output to stdout instead of lockfile")
}
