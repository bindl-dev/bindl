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

package log

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	l := log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.Disabled)

	Log = &l
}

var Log *zerolog.Logger

func IsSilent() bool {
	return Log.GetLevel() == zerolog.Disabled
}

func SetLevel(level string) error {
	lv, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		return err
	}

	l := Log.Level(lv)
	Log = &l
	Log.Debug().Str("lvl", level).Send()
	return nil
}
