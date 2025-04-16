package types

import (
	"encoding/json"

	"github.com/urfave/cli/v3"
)

// MapFlag is a flag type which requires valid JSON string maps.
type (
	MapFlag = cli.FlagBase[map[string]string, MapConfig, Map]
)

// MapConfig defines the configuration for map flags.
type MapConfig struct {
	// Any config options can be added here if needed
}

// Map implements the Value and ValueCreator interfaces for strict JSON maps.
type Map struct {
	destination *map[string]string
}

// Create implements the ValueCreator interface.
func (m Map) Create(val map[string]string, p *map[string]string, _ MapConfig) cli.Value {
	*p = val

	return &Map{
		destination: p,
	}
}

// ToString implements the ValueCreator interface.
func (m Map) ToString(val map[string]string) string {
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
func (m *Map) Set(value string) error {
	if value == "" {
		*m.destination = map[string]string{}

		return nil
	}

	return json.Unmarshal([]byte(value), m.destination)
}

// Get implements the flag.Value interface.
func (m *Map) Get() any {
	return *m.destination
}

// String implements the flag.Value interface.
func (m *Map) String() string {
	if m.destination == nil || len(*m.destination) == 0 {
		return ""
	}

	jsonBytes, err := json.Marshal(*m.destination)
	if err != nil {
		return ""
	}

	return string(jsonBytes)
}
