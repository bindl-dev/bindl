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
	"fmt"

	"github.com/spf13/cobra"

	"go.xargs.dev/bindl/command"
)

var bindlGetAll bool

var BindlGet = &cobra.Command{
	Use:   "get NAME",
	Short: "Get local copy of program",
	Long: `Get downloads the names program, which must already exist in bindl.yaml,
and ensures the program is ready to be used by setting executable flag.`,
	PreRunE: func(cmd *cobra.Command, names []string) error {
		if !bindlGetAll && len(names) == 0 {
			return fmt.Errorf("program name required but missing: specify program name or use '--all'")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, names []string) error {
		return command.LockfileProgramCommandMapper(
			cmd.Context(),
			defaultConfig,
			names,
			command.Get)
	},
}

func init() {
	BindlGet.Flags().BoolVarP(&bindlGetAll, "all", "a", false, "get all programs defined in 'bindl-lock.yaml'")
}
