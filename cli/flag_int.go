package cli

import (
	"fmt"
	"strconv"

	"github.com/urfave/cli/v3"
)

// IntFlag is a flag type which supports integer values.
// Deprecated: Use cli.IntFlag instead.
type (
	IntFlag = cli.FlagBase[int, IntConfig, Int]
)

// IntConfig defines the configuration for integer flags.
type IntConfig struct {
	// Any config options can be added here if needed
}

// Int implements the Value and ValueCreator interfaces for integers.
type Int struct {
	destination *int
}

// Create implements the ValueCreator interface.
func (i Int) Create(v int, p *int, _ IntConfig) cli.Value {
	*p = v

	return &Int{
		destination: p,
	}
}

// ToString implements the ValueCreator interface.
func (i Int) ToString(v int) string {
	return fmt.Sprintf("%d", v)
}

// Set implements the flag.Value interface.
func (i *Int) Set(v string) error {
	if v == "" {
		*i.destination = 0

		return nil
	}

	val, err := strconv.Atoi(v)
	if err != nil {
		return err
	}

	*i.destination = val

	return nil
}

// Get implements the flag.Value interface.
func (i *Int) Get() any {
	return *i.destination
}

// String implements the flag.Value interface.
func (i *Int) String() string {
	if i.destination == nil {
		return "0"
	}

	return fmt.Sprintf("%d", *i.destination)
}
