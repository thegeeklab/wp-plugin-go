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
func (s StringMap) Create(v map[string]string, p *map[string]string, _ StringMapConfig) cli.Value {
	*p = map[string]string{}

	if v != nil {
		*p = v
	}

	return &StringMap{
		destination: p,
	}
}

// ToString implements the ValueCreator interface.
func (s StringMap) ToString(v map[string]string) string {
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
func (s *StringMap) Set(v string) error {
	*s.destination = map[string]string{}

	if v == "" {
		return nil
	}

	err := json.Unmarshal([]byte(v), s.destination)
	if err != nil {
		(*s.destination)["*"] = v
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
