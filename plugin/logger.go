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
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

func loggingFlags(category string) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "log-level",
			Usage:    "plugin log level",
			EnvVars:  []string{"PLUGIN_LOG_LEVEL"},
			Value:    "info",
			Category: category,
		},
	}
}

// SetupConsoleLogger sets up the console logger.
func SetupConsoleLogger(c *cli.Context) error {
	level := "info"

	if c != nil {
		level = c.String("log-level")
	}

	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		log.Fatal().Msgf("unknown logging level: %s", level)
	}

	zerolog.SetGlobalLevel(lvl)
	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:          os.Stdout,
		PartsExclude: []string{zerolog.TimestampFieldName},
	}).With().Timestamp().Logger()

	if zerolog.GlobalLevel() <= zerolog.DebugLevel {
		log.Logger = log.With().Caller().Logger()
		log.Info().Msgf("LogLevel = %s", zerolog.GlobalLevel().String())
	}

	return nil
}
