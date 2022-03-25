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

var bindlGenerateIgnorePath = ".gitignore"

var BindlGenerateIgnore = &cobra.Command{
	Use:   "ignore",
	Short: "Generate ignore file for bindl programs",
	Long: `Generate ignore file for bindl programs 

By default, Bindl will take ".gitignore" as input and append 
<output directory>/* if it doesn't already exist.

For example, with default output directory "bin":

  $ bindl ignore -f .gitignore
  $ tail -n 1 .gitignore
  bin/*

Supports typical ignore files. e.g. .dockerignore`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return command.UpdateIgnoreFile(conf, bindlGenerateIgnorePath)
	},
}

func init() {
	BindlGenerateIgnore.Flags().StringVarP(&bindlGenerateIgnorePath, "path", "p", bindlGenerateIgnorePath, "path to ignore file")
}
