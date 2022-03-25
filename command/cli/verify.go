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

var BindlVerify = &cobra.Command{
	Use:     "verify [name, ...]",
	Aliases: []string{"validate"},
	Short:   "Verify current installation of a program",
	Long: `Verify if the currently installed program matches the specified checksum in lockfile.
If no program name is specified through args, then all programs in lockfile
will be selected.`,
	RunE: func(cmd *cobra.Command, names []string) error {
		return command.IterateLockfilePrograms(
			cmd.Context(),
			conf,
			names,
			command.Verify)
	},
}
