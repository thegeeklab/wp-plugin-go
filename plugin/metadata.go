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

// Metadata defines runtime metadata.
type Metadata struct {
	Repository Repository
	Pipeline   Pipeline
	Curr       Commit
	Prev       Commit
	Step       Step
	System     System
}

// MetadataFromContext creates a Metadata from the cli.Context.
func MetadataFromContext(ctx *cli.Context) Metadata {
	return Metadata{
		Repository: repositoryFromContext(ctx),
		Pipeline:   pipelineFromContext(ctx),
		Curr:       currFromContext(ctx),
		Prev:       prevFromContext(ctx),
		Step:       stepFromContext(ctx),
		System:     systemFromContext(ctx),
	}
}
