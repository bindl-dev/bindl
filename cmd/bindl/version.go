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

package main

import (
	"fmt"

	"github.com/spf13/cobra"
	internalversion "go.xargs.dev/bindl/internal/version"
)

var (
	version   = ""
	commit    = ""
	date      = ""
	goVersion = ""
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show build version",
	Run:   printVersion,
}

func init() {
	if version == "" {
		version = internalversion.Version()
	}
	internalversion.MarkModified(&version)
	if goVersion == "" {
		goVersion = internalversion.GoVersion()
	}
}

func printVersion(*cobra.Command, []string) {
	fmt.Printf("version: %s (%s)\ncommit: %s\ndate: %s\n", version, goVersion, commit, date)
}
