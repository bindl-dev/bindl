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
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = ""
	date    = ""
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show build version",
	Run:   printVersion,
}

func printVersion(*cobra.Command, []string) {
	fmt.Printf("version: %s\ncommit: %s\ndate: %s\n", version, commit, date)
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Printf("unable to retrieve build info")
		os.Exit(1)
	}
	fmt.Printf("go: %s\n", bi.GoVersion)
}
