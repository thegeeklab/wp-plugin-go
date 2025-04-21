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
	"strings"

	"github.com/urfave/cli/v3"
)

type (
	// Commit defines runtime metadata for a commit.
	Commit struct {
		URL          string
		SHA          string
		Ref          string
		Refspec      string
		PullRequest  int64
		SourceBranch string
		TargetBranch string
		Branch       string
		Tag          string
		Message      string
		Title        string
		Description  string
		Author       Author
	}

	// Author defines runtime metadata for a commit author.
	Author struct {
		Name   string
		Email  string
		Avatar string
	}
)

func currFlags(category string) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "commit.url",
			Usage:    "commit URL",
			Sources:  cli.EnvVars("CI_COMMIT_URL"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     "commit.sha",
			Usage:    "commit SHA",
			Sources:  cli.EnvVars("CI_COMMIT_SHA"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     "commit.ref",
			Usage:    "commit ref",
			Sources:  cli.EnvVars("CI_COMMIT_REF"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     "commit.refspec",
			Usage:    "commit refspec",
			Sources:  cli.EnvVars("CI_COMMIT_REFSPEC"),
			Category: category,
		},
		&cli.Int64Flag{
			Name:     "commit.pull-request",
			Usage:    "commit pull request",
			Sources:  cli.EnvVars("CI_COMMIT_PULL_REQUEST"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     "commit.source-branch",
			Usage:    "commit source branch",
			Sources:  cli.EnvVars("CI_COMMIT_SOURCE_BRANCH"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     "commit.target-branch",
			Usage:    "commit target branch",
			Sources:  cli.EnvVars("CI_COMMIT_TARGET_BRANCH"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     "commit.branch",
			Usage:    "commit branch",
			Sources:  cli.EnvVars("CI_COMMIT_BRANCH"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     "commit.tag",
			Usage:    "commit tag",
			Sources:  cli.EnvVars("CI_COMMIT_TAG"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     "commit.message",
			Usage:    "commit message",
			Sources:  cli.EnvVars("CI_COMMIT_MESSAGE"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     "commit.author.name",
			Usage:    "commit author name",
			Sources:  cli.EnvVars("CI_COMMIT_AUTHOR"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     "commit.author.email",
			Usage:    "commit author email",
			Sources:  cli.EnvVars("CI_COMMIT_AUTHOR_EMAIL"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     "commit.author.avatar",
			Usage:    "commit author avatar",
			Sources:  cli.EnvVars("CI_COMMIT_AUTHOR_AVATAR"),
			Category: category,
		},
	}
}

func currFromContext(c *cli.Command) Commit {
	commitTitle, commitDesc := splitMessage(c.String("commit.message"))

	return Commit{
		URL:          c.String("commit.url"),
		SHA:          c.String("commit.sha"),
		Ref:          c.String("commit.ref"),
		Refspec:      c.String("commit.refspec"),
		PullRequest:  c.Int64("commit.pull-request"),
		SourceBranch: c.String("commit.source-branch"),
		TargetBranch: c.String("commit.target-branch"),
		Branch:       c.String("commit.branch"),
		Tag:          c.String("commit.tag"),
		Message:      c.String("commit.message"),
		Title:        commitTitle,
		Description:  commitDesc,
		Author: Author{
			Name:   c.String("commit.author.name"),
			Email:  c.String("commit.author.email"),
			Avatar: c.String("commit.author.avatar"),
		},
	}
}

func prevFlags(category string) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "prev.commit.url",
			Usage:    "previous commit URL",
			Sources:  cli.EnvVars("CI_PREV_COMMIT_URL"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     "prev.commit.sha",
			Usage:    "previous commit SHA",
			Sources:  cli.EnvVars("CI_PREV_COMMIT_SHA"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     "prev.commit.ref",
			Usage:    "previous commit ref",
			Sources:  cli.EnvVars("CI_PREV_COMMIT_REF"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     "prev.commit.refspec",
			Usage:    "previous commit refspec",
			Sources:  cli.EnvVars("CI_PREV_COMMIT_REFSPEC"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     "prev.commit.branch",
			Usage:    "previous commit branch",
			Sources:  cli.EnvVars("CI_PREV_COMMIT_BRANCH"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     "prev.commit.message",
			Usage:    "previous commit message",
			Sources:  cli.EnvVars("CI_PREV_COMMIT_MESSAGE"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     "prev.commit.author.name",
			Usage:    "previous commit author name",
			Sources:  cli.EnvVars("CI_PREV_COMMIT_AUTHOR"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     "prev.commit.author.email",
			Usage:    "previous commit author email",
			Sources:  cli.EnvVars("CI_PREV_COMMIT_AUTHOR_EMAIL"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     "prev.commit.author.avatar",
			Usage:    "previous commit author avatar",
			Sources:  cli.EnvVars("CI_PREV_COMMIT_AUTHOR_AVATAR"),
			Category: category,
		},
	}
}

func prevFromContext(c *cli.Command) Commit {
	commitTitle, commitDesc := splitMessage(c.String("commit.message"))

	return Commit{
		URL:         c.String("prev.commit.url"),
		SHA:         c.String("prev.commit.sha"),
		Ref:         c.String("prev.commit.ref"),
		Refspec:     c.String("prev.commit.refspec"),
		Branch:      c.String("prev.commit.branch"),
		Message:     c.String("prev.commit.message"),
		Title:       commitTitle,
		Description: commitDesc,
		Author: Author{
			Name:   c.String("prev.commit.author.name"),
			Email:  c.String("prev.commit.author.email"),
			Avatar: c.String("prev.commit.author.avatar"),
		},
	}
}

// splitMessage splits a commit message into a title and description.
// It splits the message on the first newline character, with the first
// line as the title, and the rest as the description. If there is no newline,
// the entire message is returned as the title, and the description is empty.
func splitMessage(message string) (string, string) {
	//nolint:mnd
	switch parts := strings.SplitN(message, "\n", 2); len(parts) {
	case 1:
		return parts[0], ""
	//nolint:mnd
	case 2:
		return parts[0], parts[1]
	}

	return "", ""
}
