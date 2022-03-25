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
	"github.com/bindl-dev/bindl/command"
	"github.com/spf13/cobra"
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
		return command.Sync(cmd.Context(), conf, bindlSyncStdout)
	},
}

func init() {
	BindlSync.Flags().BoolVar(&bindlSyncStdout, "stdout", false, "write output to stdout instead of lockfile")
}
