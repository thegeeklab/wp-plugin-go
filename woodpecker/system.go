// Copyright 2023 Woodpecker Authors
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

package plugin

import (
	"github.com/urfave/cli/v2"
)

// System defines runtime metadata for a ci/cd system.
type System struct {
	Name     string `json:"name,omitempty"`
	Host     string `json:"host,omitempty"`
	Link     string `json:"link,omitempty"`
	Platform string `json:"arch,omitempty"`
	Version  string `json:"version,omitempty"`
}

func systemFlags(category string) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "system.name",
			Usage:    "system name",
			EnvVars:  []string{"CI_SYSTEM_NAME"},
			Category: category,
		},
		&cli.StringFlag{
			Name:     "system.host",
			Usage:    "system host",
			EnvVars:  []string{"CI_SYSTEM_HOST"},
			Category: category,
		},
		&cli.StringFlag{
			Name:     "system.link",
			Usage:    "system link",
			EnvVars:  []string{"CI_SYSTEM_LINK"},
			Category: category,
		},
		&cli.StringFlag{
			Name:     "system.arch",
			Usage:    "system arch",
			EnvVars:  []string{"CI_SYSTEM_PLATFORM"},
			Category: category,
		},
		&cli.StringFlag{
			Name:     "system.version",
			Usage:    "system version",
			EnvVars:  []string{"CI_SYSTEM_VERSION"},
			Category: category,
		},
	}
}

func systemFromContext(ctx *cli.Context) System {
	link := ctx.String("system.link")
	host := ctx.String("system.host")

	return System{
		Name:     ctx.String("system.name"),
		Host:     host,
		Link:     link,
		Platform: ctx.String("system.arch"),
		Version:  ctx.String("system.version"),
	}
}
