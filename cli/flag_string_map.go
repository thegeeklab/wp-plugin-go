package cli

import (
	"encoding/json"

	"github.com/urfave/cli/v3"
)

// StringMapFlag is a flag type which supports JSON string maps.
type (
	StringMapFlag = cli.FlagBase[map[string]string, StringMapConfig, StringMap]
)

// StringMapConfig defines the configuration for string map flags.
type StringMapConfig struct {
	// Any config options can be added here if needed
}

// StringMap implements the Value and ValueCreator interfaces for string maps.
type StringMap struct {
	destination *map[string]string
}

// Create implements the ValueCreator interface.
func (s StringMap) Create(val map[string]string, p *map[string]string, _ StringMapConfig) cli.Value {
	*p = val

	return &StringMap{
		destination: p,
	}
}

// ToString implements the ValueCreator interface.
func (s StringMap) ToString(val map[string]string) string {
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
func (s *StringMap) Set(value string) error {
	if value == "" {
		*s.destination = map[string]string{}

		return nil
	}

	err := json.Unmarshal([]byte(value), s.destination)
	if err != nil {
		// Initialize the map if it's nil
		if *s.destination == nil {
			*s.destination = map[string]string{}
		}

		// For StringMapFlag, we want to handle non-JSON values as a wildcard key
		(*s.destination)["*"] = value
	}

	return nil
}

// Get implements the flag.Value interface.
func (s *StringMap) Get() any {
	return *s.destination
}

// String implements the flag.Value interface.
func (s *StringMap) String() string {
	if s.destination == nil || len(*s.destination) == 0 {
		return ""
	}

	jsonBytes, err := json.Marshal(*s.destination)
	if err != nil {
		return ""
	}

	return string(jsonBytes)
}
