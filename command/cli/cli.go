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
	"github.com/spf13/cobra"
	"go.xargs.dev/bindl/config"
	"go.xargs.dev/bindl/internal/log"
)

var All = []*cobra.Command{
	BindlGet,
	BindlSync,
	BindlGenerate,
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
