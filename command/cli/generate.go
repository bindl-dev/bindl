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

import "github.com/spf13/cobra"

var BindlGenerate = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "Generate project integration files (e.g. Makefile, .gitignore)",
	Long:    `Generate common setup for projects to have smooth workflow.`,
}

func init() {
	BindlGenerate.AddCommand(BindlGenerateIgnore)
	BindlGenerate.AddCommand(BindlGenerateMake)
}
