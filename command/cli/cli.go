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

	"github.com/bindl-dev/bindl/config"
	"github.com/bindl-dev/bindl/internal/log"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
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
		if err := envconfig.Process("BINDL", conf); err != nil {
			return err
		}

		var logLevel string
		switch {
		case conf.Silent:
			logLevel = "disabled"
		case conf.Debug:
			logLevel = "debug"
		default:
			logLevel = "info"
		}
		return log.SetLevel(logLevel)
	},
}

var conf = &config.Runtime{
	Path:         "./bindl.yaml",
	LockfilePath: "./.bindl-lock.yaml",
	BinDir:       "./bin",
	ProgDir:      ".bindl/programs",

	UseCache: false,

	Debug:  false,
	Silent: false,

	OS:   runtime.GOOS,
	Arch: runtime.GOARCH,
}

func init() {
	Root.PersistentFlags().StringVarP(&conf.Path, "config", "c", conf.Path, "path to configuration file")
	Root.PersistentFlags().StringVarP(&conf.LockfilePath, "lock", "l", conf.LockfilePath, "path to lockfile")
	Root.PersistentFlags().StringVarP(&conf.BinDir, "bin", "b", conf.BinDir, "directory in PATH to add binaries")
	Root.PersistentFlags().StringVar(&conf.ProgDir, "prog", conf.ProgDir, "directory to save real binary content")

	Root.PersistentFlags().BoolVar(&conf.UseCache, "cache", conf.UseCache, "read and write cache")

	Root.PersistentFlags().BoolVarP(&conf.Silent, "silent", "s", conf.Silent, "silence logs")
	Root.PersistentFlags().BoolVar(&conf.Debug, "debug", conf.Debug, "show debug logs")
}
