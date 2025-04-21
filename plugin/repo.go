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
	"github.com/urfave/cli/v3"
)

// Repository defines runtime metadata for a repository.
type Repository struct {
	Slug     string
	Name     string
	Owner    string
	URL      string
	CloneURL string
	Private  bool
	Branch   string
	RemoteID int64
}

func repositoryFlags(category string) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "repo.slug",
			Usage:    "repo slug",
			Sources:  cli.EnvVars("CI_REPO"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     "repo.name",
			Usage:    "repo name",
			Sources:  cli.EnvVars("CI_REPO_NAME"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     "repo.owner",
			Usage:    "repo owner",
			Sources:  cli.EnvVars("CI_REPO_OWNER"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     "repo.url",
			Usage:    "repo url",
			Sources:  cli.EnvVars("CI_REPO_URL"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     "repo.clone-url",
			Usage:    "repo clone url",
			Sources:  cli.EnvVars("CI_REPO_CLONE_URL"),
			Category: category,
		},
		&cli.BoolFlag{
			Name:     "repo.private",
			Usage:    "repo private",
			Sources:  cli.EnvVars("CI_REPO_PRIVATE"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     "repo.default-branch",
			Usage:    "repo default branch",
			Sources:  cli.EnvVars("CI_REPO_DEFAULT_BRANCH"),
			Category: category,
		},
		&cli.Int64Flag{
			Name:     "repo.remote-id",
			Usage:    "repo remote id",
			Sources:  cli.EnvVars("CI_REPO_REMOTE_ID"),
			Category: category,
		},
	}
}

func repositoryFromContext(c *cli.Command) Repository {
	return Repository{
		Slug:     c.String("repo.slug"),
		Name:     c.String("repo.name"),
		Owner:    c.String("repo.owner"),
		URL:      c.String("repo.url"),
		CloneURL: c.String("repo.clone-url"),
		Private:  c.Bool("repo.private"),
		Branch:   c.String("repo.default-branch"),
		RemoteID: c.Int64("repo.remote-id"),
	}
}
