package plugin

import (
	"errors"
	"fmt"
	"os"

	"github.com/thegeeklab/wp-plugin-go/v4/types"
	"github.com/urfave/cli/v2"
)

var ErrTypeAssertionFailed = errors.New("type assertion failed")

type Environment map[string]string

// Lookup retrieves the value for the given key from the Environment. If the key
// is not found in the Environment, it falls back to looking up the value in
// the OS environment.
func (e Environment) Lookup(key string) (string, bool) {
	value, ok := e[key]
	if ok {
		return value, ok
	}

	return os.LookupEnv(key)
}

// Value returns a slice of strings representing the key-value pairs of the
// Environment. Each string is formatted as "key=value".
func (e Environment) Value() []string {
	values := make([]string, 0, len(e))
	for key, value := range e {
		values = append(values, fmt.Sprintf("%s=%s", key, value))
	}
	return values
}

func environmentFlags(category string) []cli.Flag {
	return []cli.Flag{
		&cli.GenericFlag{
			Name:     "environment",
			Usage:    "plugin environment variables",
			EnvVars:  []string{"PLUGIN_ENVIRONMENT"},
			Value:    &types.StringMapFlag{},
			Category: category,
		},
	}
}

func EnvironmentFromContext(ctx *cli.Context) (Environment, error) {
	env, ok := ctx.Generic("env").(*types.StringMapFlag)
	if !ok {
		return nil, fmt.Errorf("%w: failed to read plugin environment", ErrTypeAssertionFailed)
	}

	return env.Get(), nil
}
