package cli

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"
)

var errParseDeepStringMap = errors.New("unable to parse flag value as JSON string map")

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
	v map[string]map[string]string,
	p *map[string]map[string]string,
	_ DeepStringMapConfig,
) cli.Value {
	*p = map[string]map[string]string{}

	if v != nil {
		*p = v
	}

	return &DeepStringMap{
		destination: p,
	}
}

// ToString implements the ValueCreator interface.
func (d DeepStringMap) ToString(v map[string]map[string]string) string {
	if len(v) == 0 {
		return ""
	}

	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return ""
	}

	return string(jsonBytes)
}

// Set implements the flag.Value interface.
// Valid nested JSON objects are parsed directly. JSON values that are not
// strings (numbers, booleans, null) are coerced to their string
// representations. Flat (non-nested) JSON objects are stored under the "*"
// key with the same coercion applied. Input that is not valid JSON returns
// an error.
func (d *DeepStringMap) Set(v string) error {
	*d.destination = map[string]map[string]string{}

	if v == "" {
		return nil
	}

	if err := json.Unmarshal([]byte(v), d.destination); err == nil {
		return nil
	}

	var rawNested map[string]map[string]any
	if json.Unmarshal([]byte(v), &rawNested) == nil {
		*d.destination = convertNestedMap(rawNested)

		return nil
	}

	var rawSingle map[string]any
	if json.Unmarshal([]byte(v), &rawSingle) == nil {
		*d.destination = map[string]map[string]string{
			"*": convertFlatMap(rawSingle),
		}

		return nil
	}

	return fmt.Errorf("%w: %q", errParseDeepStringMap, v)
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
