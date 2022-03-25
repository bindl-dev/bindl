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
	"go.xargs.dev/bindl/config"
)

var bindlListOneline = false

var BindlList = &cobra.Command{
	Use:   "list",
	Short: "List out programs defined in lockfile",
	Long:  "List shows all program names defined in lockfile.",

	Args: cobra.NoArgs,

	RunE: func(cmd *cobra.Command, args []string) error {
		l, err := config.ParseLock(conf.LockfilePath)
		if err != nil {
			return fmt.Errorf("parsing lockfile: %w", err)
		}
		separator := "\n"
		if bindlListOneline {
			separator = " "
		}
		for _, p := range l.Programs {
			fmt.Printf("%s%s", p.Name, separator)
		}
		if bindlListOneline {
			fmt.Println()
		}
		return nil
	},
}

func init() {
	BindlList.Flags().BoolVar(&bindlListOneline, "oneline", bindlListOneline, "List programs in one line, space separated")
}
