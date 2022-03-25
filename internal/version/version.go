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

package version

import (
	"runtime/debug"
)

var (
	version   string
	commit    string
	date      string
	goVersion string
	modified  bool
)

func init() {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		panic("unable to retrieve build info")
	}
	version = bi.Main.Version
	goVersion = bi.GoVersion

	for _, s := range bi.Settings {
		switch s.Key {
		case "vcs.revision":
			commit = s.Value
		case "vcs.time":
			date = s.Value
		case "vcs.modified":
			modified = s.Value == "true"
			version = version + "-dirty"
		default:
			continue
		}
	}
}

func Version() string {
	return version
}

func Commit() string {
	return commit
}

func Date() string {
	return date
}

func GoVersion() string {
	return goVersion
}

func Modified() bool {
	return modified
}

func MarkModified(v *string) {
	if Modified() {
		*v = *v + "-dirty"
	}
}

func SetFromCmd(v string) {
	version = v
}

func Summary() string {
	return version + " (" + commit + ") (" + goVersion + ")"
}
