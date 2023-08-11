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

	"github.com/urfave/cli/v2"
)

// Pipeline defines runtime metadata for a pipeline.
type Pipeline struct {
	Number       int64     `json:"number,omitempty"`
	Status       string    `json:"status,omitempty"`
	Event        string    `json:"event,omitempty"`
	Link         string    `json:"link,omitempty"`
	DeployTarget string    `json:"target,omitempty"`
	Created      time.Time `json:"created,omitempty"`
	Started      time.Time `json:"started,omitempty"`
	Finished     time.Time `json:"finished,omitempty"`
	Parent       int64     `json:"parent,omitempty"`
}

func pipelineFlags(category string) []cli.Flag {
	return []cli.Flag{
		&cli.Int64Flag{
			Name:     "pipeline.number",
			Usage:    "pipeline number",
			EnvVars:  []string{"CI_PIPELINE_NUMBER"},
			Category: category,
		},
		&cli.StringFlag{
			Name:     "pipeline.status",
			Usage:    "pipeline status",
			EnvVars:  []string{"CI_PIPELINE_STATUS"},
			Category: category,
		},
		&cli.StringFlag{
			Name:     "pipeline.event",
			Usage:    "pipeline event",
			EnvVars:  []string{"CI_PIPELINE_EVENT"},
			Category: category,
		},
		&cli.StringFlag{
			Name:     "pipeline.link",
			Usage:    "pipeline link",
			EnvVars:  []string{"CI_PIPELINE_LINK"},
			Category: category,
		},
		&cli.StringFlag{
			Name:     "pipeline.deploy-target",
			Usage:    "pipeline deployment target",
			EnvVars:  []string{"CI_PIPELINE_DEPLOY_TARGET"},
			Category: category,
		},
		&cli.Int64Flag{
			Name:     "pipeline.created",
			Usage:    "pipeline creation time",
			EnvVars:  []string{"CI_PIPELINE_CREATED"},
			Category: category,
		},
		&cli.Int64Flag{
			Name:     "pipeline.started",
			Usage:    "pipeline start time",
			EnvVars:  []string{"CI_PIPELINE_STARTED"},
			Category: category,
		},
		&cli.Int64Flag{
			Name:     "pipeline.finished",
			Usage:    "pipeline finish time",
			EnvVars:  []string{"CI_PIPELINE_FINISHED"},
			Category: category,
		},
		&cli.Int64Flag{
			Name:     "pipeline.parent",
			Usage:    "pipeline parent",
			EnvVars:  []string{"CI_PIPELINE_PARENT"},
			Category: category,
		},
	}
}

func pipelineFromContext(c *cli.Context) Pipeline {
	return Pipeline{
		Number:       c.Int64("pipeline.number"),
		Status:       c.String("pipeline.status"),
		Event:        c.String("pipeline.event"),
		Link:         c.String("pipeline.link"),
		DeployTarget: c.String("pipeline.deploy-target"),
		Created:      time.Unix(c.Int64("pipeline.created"), 0),
		Started:      time.Unix(c.Int64("pipeline.started"), 0),
		Finished:     time.Unix(c.Int64("pipeline.finished"), 0),
		Parent:       c.Int64("pipeline.parent"),
	}
}
