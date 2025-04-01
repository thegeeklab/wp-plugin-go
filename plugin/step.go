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
	"time"

	"github.com/urfave/cli/v3"
)

// Step defines runtime metadata for a step.
type Step struct {
	Number   int
	Started  time.Time
	Finished time.Time
}

func stepFlags(category string) []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:     "step.number",
			Usage:    "step number",
			EnvVars:  []string{"CI_STEP_NUMBER"},
			Category: category,
		},
		&cli.Int64Flag{
			Name:     "step.started",
			Usage:    "step start time",
			EnvVars:  []string{"CI_STEP_STARTED"},
			Category: category,
		},
		&cli.Int64Flag{
			Name:     "step.finished",
			Usage:    "step finish time",
			EnvVars:  []string{"CI_STEP_FINISHED"},
			Category: category,
		},
	}
}

func stepFromContext(c *cli.Context) Step {
	return Step{
		Number:   c.Int("step.number"),
		Started:  time.Unix(c.Int64("step.started"), 0),
		Finished: time.Unix(c.Int64("step.finished"), 0),
	}
}
