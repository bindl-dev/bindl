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
	"go.xargs.dev/bindl/command"
)

var bindlPurgeAll = false
var bindlPurgeDryRun = false

var BindlPurge = &cobra.Command{
	Use:   "purge",
	Short: "Remove downloaded programs",
	Long: `Remove downloaded programs from cache, which are not listed in the lockfile.
Passing --all would remove all existing programs regardless of lockfile.`,
	RunE: func(cmd *cobra.Command, names []string) error {
		return command.Purge(cmd.Context(), defaultConfig, bindlPurgeAll, bindlPurgeDryRun)
	},
}

func init() {
	BindlPurge.Flags().BoolVarP(&bindlPurgeAll, "all", "a", bindlPurgeAll, "purge all existing programs, regardless of lockfile")
	BindlPurge.Flags().BoolVar(&bindlPurgeDryRun, "dry-run", bindlPurgeDryRun, "dry-run purge (print paths without deleting)")
}
