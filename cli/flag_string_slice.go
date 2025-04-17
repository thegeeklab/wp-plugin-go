package cli

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"
)

// StringSliceFlag is a flag type which support comma separated values and escaping to not split at unwanted lines.
type (
	StringSliceFlag = cli.FlagBase[[]string, StringSliceConfig, StringSlice]
)

// StringConfig defines the configuration for string flags.
type StringSliceConfig struct {
	Separator    string
	EscapeString string
}

// StringSlice implements the Value and ValueCreator interfaces for string slices.
type StringSlice struct {
	destination  *[]string
	separator    string
	escapeString string
}

// Create implements the ValueCreator interface.
func (s StringSlice) Create(val []string, p *[]string, c StringSliceConfig) cli.Value {
	*p = val

	return &StringSlice{
		destination:  p,
		separator:    c.Separator,
		escapeString: c.EscapeString,
	}
}

// ToString implements the ValueCreator interface.
func (s StringSlice) ToString(val []string) string {
	if len(val) == 0 {
		return ""
	}

	return fmt.Sprintf("%q", strings.Join(val, s.separator))
}

// Set implements the flag.Value interface.
func (s *StringSlice) Set(value string) error {
	if value == "" {
		*s.destination = []string{}

		return nil
	}

	out := strings.Split(value, s.separator)

	//nolint:mnd
	for i := len(out) - 2; i >= 0; i-- {
		if strings.HasSuffix(out[i], s.escapeString) {
			out[i] = out[i][:len(out[i])-len(s.escapeString)] + s.separator + out[i+1]
			out = append(out[:i+1], out[i+2:]...)
		}
	}

	*s.destination = out

	return nil
}

// Get implements the flag.Value interface.
func (s *StringSlice) Get() any {
	return *s.destination
}

// String implements the flag.Value interface.
func (s *StringSlice) String() string {
	if s.destination == nil || len(*s.destination) == 0 {
		return ""
	}

	return strings.Join(*s.destination, s.separator)
}
