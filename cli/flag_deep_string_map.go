package cli

import (
	"encoding/json"

	"github.com/urfave/cli/v3"
)

// DeepStringMapFlag is a flag type which supports nested JSON string maps.
type (
	DeepStringMapFlag = cli.FlagBase[map[string]map[string]string, DeepStringMapConfig, DeepStringMap]
)

// DeepStringMapConfig defines the configuration for deep string map flags.
type DeepStringMapConfig struct {
	// Any config options can be added here if needed
}

// DeepStringMap implements the Value and ValueCreator interfaces for nested string maps.
type DeepStringMap struct {
	destination *map[string]map[string]string
}

// Create implements the ValueCreator interface.
func (d DeepStringMap) Create(
	val map[string]map[string]string,
	p *map[string]map[string]string,
	_ DeepStringMapConfig,
) cli.Value {
	*p = val

	return &DeepStringMap{
		destination: p,
	}
}

// ToString implements the ValueCreator interface.
func (d DeepStringMap) ToString(val map[string]map[string]string) string {
	if len(val) == 0 {
		return ""
	}

	jsonBytes, err := json.Marshal(val)
	if err != nil {
		return ""
	}

	return string(jsonBytes)
}

// Set implements the flag.Value interface.
func (d *DeepStringMap) Set(value string) error {
	if value == "" {
		*d.destination = map[string]map[string]string{}

		return nil
	}

	err := json.Unmarshal([]byte(value), d.destination)
	if err != nil {
		// Try to parse as a single-level map
		single := map[string]string{}

		err := json.Unmarshal([]byte(value), &single)
		if err != nil {
			return err
		}

		// Store the single-level map under a wildcard key
		(*d.destination) = map[string]map[string]string{
			"*": single,
		}
	}

	return nil
}

// Get implements the flag.Value interface.
func (d *DeepStringMap) Get() any {
	return *d.destination
}

// String implements the flag.Value interface.
func (d *DeepStringMap) String() string {
	if d.destination == nil || len(*d.destination) == 0 {
		return ""
	}

	jsonBytes, err := json.Marshal(*d.destination)
	if err != nil {
		return ""
	}

	return string(jsonBytes)
}
