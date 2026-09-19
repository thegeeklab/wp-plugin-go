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
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

//nolint:lll
const appHelpTemplate = `NAME:
   {{template "helpNameTemplate" .}}

USAGE:
   {{if .UsageText}}{{wrap .UsageText 3}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}{{if .Args}}[arguments...]{{end}}{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}

VERSION:
   {{.Version}}{{end}}{{end}}{{if .Description}}

DESCRIPTION:
   {{template "descriptionTemplate" .}}{{end}}
{{- if len .Authors}}

AUTHOR{{template "authorsTemplate" .}}{{end}}{{if .VisibleCommands}}

COMMANDS:{{template "visibleCommandCategoryTemplate" .}}{{end}}{{if .VisibleFlagCategories}}

GLOBAL OPTIONS:{{range .VisibleFlagCategories}}{{if and .Name (ne .Name "Plugin Flags")}}{{continue}}{{end}}
   {{if .Name}}{{.Name}}

   {{end}}{{$flglen := len .Flags}}{{range $i, $e := .Flags}}{{if eq (subtract $flglen $i) 1}}{{$e}}
   {{else}}{{$e}}
   {{end}}{{end}}{{end}}{{else if .VisibleFlags}}

GLOBAL OPTIONS:{{template "visibleFlagTemplate" .}}{{end}}{{if .Copyright}}

COPYRIGHT:
   {{template "copyrightTemplate" .}}{{end}}
`

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
	// Hide woodpecker system flags.
	HideWoodpeckerFlags bool
}

// Plugin defines the plugin instance.
type Plugin struct {
	App         *cli.Command
	execute     ExecuteFunc
	network     *Network
	metadata    *Metadata
	environment Environment
}

var (
	errMetadataNotAvailable      = errors.New("metadata not available")
	errNetworkFlagsNotRegistered = errors.New("network flags not registered: add NetworkFlags to Options.Flags")
	errEnvFlagsNotRegistered     = errors.New("environment flags not registered: add EnvironmentFlags to Options.Flags")
)

// GetMetadata returns the pipeline metadata.
// Metadata is always available after the plugin action runs.
func (p *Plugin) GetMetadata() (Metadata, error) {
	if p.metadata == nil {
		return Metadata{}, errMetadataNotAvailable
	}

	return *p.metadata, nil
}

// GetNetwork returns the network configuration.
// Returns an error if NetworkFlags were not registered via Options.Flags.
func (p *Plugin) GetNetwork() (Network, error) {
	if p.network == nil {
		return Network{}, errNetworkFlagsNotRegistered
	}

	return *p.network, nil
}

// GetEnvironment returns the environment variables.
// Returns an error if EnvironmentFlags were not registered via Options.Flags.
func (p *Plugin) GetEnvironment() (Environment, error) {
	if p.environment == nil {
		return nil, errEnvFlagsNotRegistered
	}

	return p.environment, nil
}

// ExecuteFunc defines the function that is executed by the plugin.
type ExecuteFunc func(ctx context.Context) error

// New plugin instance.
func New(opt Options) *Plugin {
	if _, err := os.Stat("/run/woodpecker/env"); err == nil {
		_ = godotenv.Overload("/run/woodpecker/env")
	}

	_, _ = SetupConsoleLogger(context.Background(), nil)

	app := &cli.Command{
		Name:    opt.Name,
		Usage:   opt.Description,
		Version: opt.Version,
		Flags:   append(opt.Flags, Flags()...),
		Before:  SetupConsoleLogger,
	}

	if opt.HideWoodpeckerFlags {
		app.CustomRootCommandHelpTemplate = appHelpTemplate
	}

	cli.VersionPrinter = func(cmd *cli.Command) {
		version := fmt.Sprintf("%s version=%s %s\n", cmd.Name, cmd.Version, opt.VersionMetadata)
		fmt.Println(strings.TrimSpace(version))
	}

	plugin := &Plugin{
		App:     app,
		execute: opt.Execute,
	}
	plugin.App.Action = plugin.action

	return plugin
}

func (p *Plugin) action(ctx context.Context, cmd *cli.Command) error {
	var err error

	metadata := MetadataFromContext(cmd)
	p.metadata = &metadata

	if cmd.Value("transport.insecure-skip-verify") != nil {
		network := NetworkFromContext(cmd)
		p.network = &network
	}

	if cmd.Value("environment") != nil {
		p.environment, err = EnvironmentFromContext(cmd)
		if err != nil {
			return err
		}
	}

	if p.metadata.Pipeline.URL == "" {
		url, err := url.JoinPath(
			p.metadata.System.URL,
			"repos",
			p.metadata.Repository.Slug,
			"pipeline",
			strconv.FormatInt(p.metadata.Pipeline.Number, 10),
		)
		if err == nil {
			p.metadata.Pipeline.URL = url
		}
	}

	if p.execute == nil {
		panic("plugin execute function is not set")
	}

	return p.execute(ctx)
}

// Run the plugin.
func (p *Plugin) Run() {
	if err := p.App.Run(context.Background(), os.Args); err != nil {
		log.Error().Err(err).Msg("execution failed")
		os.Exit(1)
	}
}
