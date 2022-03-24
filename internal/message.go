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

package internal

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"go.xargs.dev/bindl/internal/log"
)

var errHeader = color.HiRedString("ERROR")

func ErrorMsg(err error) {
	if !log.IsSilent() {
		fmt.Fprintf(os.Stderr, "%s - %s\n", errHeader, err.Error())
	}
}

func Msgf(msg string, vars ...any) {
	if !log.IsSilent() {
		fmt.Fprintf(os.Stderr, msg, vars...)
	}
}
