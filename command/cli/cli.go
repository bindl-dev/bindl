// Copyright 2022 Bindl Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"runtime"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
	"go.xargs.dev/bindl/config"
	"go.xargs.dev/bindl/internal/log"
)

var All = []*cobra.Command{
	BindlGet,
	BindlSync,
	BindlList,
	BindlGenerate,
	BindlPurge,
	BindlVerify,
}

var Root = &cobra.Command{
	Use:  "bindl",
	Long: "Bindl is a static binary downloader for project development and infrastructure.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := envconfig.Process("BINDL", defaultConfig); err != nil {
			return err
		}

		var logLevel string
		switch {
		case defaultConfig.Silent:
			logLevel = "disabled"
		case defaultConfig.Debug:
			logLevel = "debug"
		default:
			logLevel = "info"
		}
		return log.SetLevel(logLevel)
	},
}

var defaultConfig = &config.Runtime{
	Path:         "./bindl.yaml",
	LockfilePath: "./.bindl-lock.yaml",
	BinDir:       "./bin",
	ProgDir:      ".bindl/programs",

	Debug:  false,
	Silent: false,

	OS:   runtime.GOOS,
	Arch: runtime.GOARCH,
}

func init() {
	Root.PersistentFlags().StringVarP(&defaultConfig.Path, "config", "c", defaultConfig.Path, "path to configuration file")
	Root.PersistentFlags().StringVarP(&defaultConfig.LockfilePath, "lock", "l", defaultConfig.LockfilePath, "path to lockfile")
	Root.PersistentFlags().StringVarP(&defaultConfig.BinDir, "bin", "b", defaultConfig.BinDir, "directory in PATH to add binaries")
	Root.PersistentFlags().StringVar(&defaultConfig.ProgDir, "prog", defaultConfig.ProgDir, "directory to save real binary content")

	Root.PersistentFlags().BoolVarP(&defaultConfig.Silent, "silent", "s", defaultConfig.Silent, "silence logs")
	Root.PersistentFlags().BoolVar(&defaultConfig.Debug, "debug", defaultConfig.Debug, "show debug logs")
}
