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

var bindlGenerateMakefilePath = "Makefile.bindl"

var BindlGenerateMake = &cobra.Command{
	Use:     "make",
	Aliases: []string{"makefile"},
	Short:   "Generate Makefile for bindl programs",
	Long: `Generate Makefile for all programs in lockfile.

By default, the generated Makefile will be named 'Makefile.bindl', which can be
imported by the project's primary Makefile using 'include' directive.

After including, you can use rules defined in 'Makefile.bindl' as a dependency
in your other rules. For example:

	$ head -n 5 Makefile
	include Makefile.bindl

	.PHONY: container
	container: bin/ko
		bin/ko publish -B .

Calling the imported rules also works on 'make' CLI.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return command.GenerateMakefile(conf, bindlGenerateMakefilePath)
	},
}

func init() {
	BindlGenerateMake.Flags().StringVarP(&bindlGenerateMakefilePath, "path", "p", bindlGenerateMakefilePath, "path to generated Makefile")
}
