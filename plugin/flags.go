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

const (
	FlagsRepositoryCategory = "Woodpecker Repository Flags"
	FlagsPipelineCategory   = "Woodpecker Pipeline Flags"
	FlagsCommitCategory     = "Woodpecker Commit Flags"
	FlagsStepCategory       = "Woodpecker Step Flags"
	FlagsSystemCategory     = "Woodpecker System Flags"
	FlagsPluginCategory     = "Plugin Flags"
)

// Flags has the cli.Flags for the Woodpecker plugin.
func Flags() []cli.Flag {
	flags := make([]cli.Flag, 0)

	// Pipeline flags
	flags = append(flags, repositoryFlags(FlagsRepositoryCategory)...)
	flags = append(flags, pipelineFlags(FlagsPipelineCategory)...)
	flags = append(flags, currFlags(FlagsCommitCategory)...)
	flags = append(flags, prevFlags(FlagsCommitCategory)...)
	flags = append(flags, stepFlags(FlagsStepCategory)...)
	flags = append(flags, systemFlags(FlagsSystemCategory)...)

	// Plugin flags
	flags = append(flags, loggingFlags(FlagsPluginCategory)...)
	flags = append(flags, networkFlags(FlagsPluginCategory)...)
	flags = append(flags, environmentFlags(FlagsPluginCategory)...)

	return flags
}
