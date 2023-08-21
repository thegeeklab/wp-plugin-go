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
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

// Options defines the options for the plugin.
type Options struct {
	// Name of the plugin.
	Name string
	// Description of the plugin.
	Description string
	// Version of the plugin.
	Version string
	// Version metadata of the plugin.
	VersionMetadata string
	// Flags of the plugin.
	Flags []cli.Flag
	// Execute function of the plugin.
	Execute ExecuteFunc
}

// Plugin defines the plugin instance.
type Plugin struct {
	App     *cli.App
	execute ExecuteFunc
	Context *cli.Context
	// Network options.
	Network Network
	// Metadata of the current pipeline.
	Metadata Metadata
}

// ExecuteFunc defines the function that is executed by the plugin.
type ExecuteFunc func(ctx context.Context) error

// New plugin instance.
func New(opt Options) *Plugin {
	if _, err := os.Stat("/run/woodpecker/env"); err == nil {
		_ = godotenv.Overload("/run/woodpecker/env")
	}

	app := &cli.App{
		Name:    opt.Name,
		Usage:   opt.Description,
		Version: opt.Version,
		Flags:   append(opt.Flags, Flags()...),
	}

	cli.VersionPrinter = func(c *cli.Context) {
		version := fmt.Sprintf("%s version=%s %s\n", c.App.Name, c.App.Version, opt.VersionMetadata)
		fmt.Print(strings.TrimSpace(version))
	}

	plugin := &Plugin{
		App:     app,
		execute: opt.Execute,
	}
	plugin.App.Action = plugin.action

	return plugin
}

func (p *Plugin) action(ctx *cli.Context) error {
	if err := SetupConsoleLogger(ctx); err != nil {
		return err
	}

	p.Metadata = MetadataFromContext(ctx)
	p.Network = NetworkFromContext(ctx)
	p.Context = ctx

	if p.Metadata.Pipeline.URL == "" {
		url, err := url.JoinPath(
			p.Metadata.System.URL,
			"repos",
			p.Metadata.Repository.Slug,
			"pipeline",
			strconv.FormatInt(p.Metadata.Pipeline.Number, 10),
		)
		if err == nil {
			p.Metadata.Pipeline.URL = url
		}
	}

	if p.execute == nil {
		panic("plugin execute function is not set")
	}

	return p.execute(ctx.Context)
}

// Run the plugin.
func (p *Plugin) Run() {
	if err := p.App.Run(os.Args); err != nil {
		log.Error().Err(err).Msg("execution failed")
		os.Exit(1)
	}
}
