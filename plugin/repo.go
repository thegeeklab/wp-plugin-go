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

// Repository defines runtime metadata for a repository.
type Repository struct {
	Name     string `json:"name,omitempty"`
	Owner    string `json:"owner,omitempty"`
	Link     string `json:"link,omitempty"`
	CloneURL string `json:"clone_url,omitempty"`
	Private  bool   `json:"private,omitempty"`
	Branch   string `json:"default_branch,omitempty"`
}

func repositoryFlags(category string) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "repo.name",
			Usage:    "repo name",
			EnvVars:  []string{"CI_REPO_NAME"},
			Category: category,
		},
		&cli.StringFlag{
			Name:     "repo.owner",
			Usage:    "repo owner",
			EnvVars:  []string{"CI_REPO_OWNER"},
			Category: category,
		},
		&cli.StringFlag{
			Name:     "repo.link",
			Usage:    "repo link",
			EnvVars:  []string{"CI_REPO_LINK"},
			Category: category,
		},
		&cli.StringFlag{
			Name:     "repo.clone-url",
			Usage:    "repo clone url",
			EnvVars:  []string{"CI_REPO_CLONE_URL"},
			Category: category,
		},
		&cli.BoolFlag{
			Name:     "repo.private",
			Usage:    "repo private",
			EnvVars:  []string{"CI_REPO_PRIVATE"},
			Category: category,
		},
		&cli.StringFlag{
			Name:     "repo.default-branch",
			Usage:    "repo default branch",
			EnvVars:  []string{"CI_REPO_DEFAULT_BRANCH"},
			Category: category,
		},
	}
}

func repositoryFromContext(c *cli.Context) Repository {
	return Repository{
		Name:     c.String("repo.name"),
		Owner:    c.String("repo.owner"),
		Link:     c.String("repo.link"),
		CloneURL: c.String("repo.clone-url"),
		Private:  c.Bool("repo.private"),
		Branch:   c.String("repo.default-branch"),
	}
}
