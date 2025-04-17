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
	Delimiter    string
	EscapeString string
}

// StringSlice implements the Value and ValueCreator interfaces for string slices.
type StringSlice struct {
	destination  *[]string
	delimiter    string
	escapeString string
}

// Create implements the ValueCreator interface.
func (s StringSlice) Create(v []string, p *[]string, c StringSliceConfig) cli.Value {
	*p = v

	return &StringSlice{
		destination:  p,
		delimiter:    c.Delimiter,
		escapeString: c.EscapeString,
	}
}

// ToString implements the ValueCreator interface.
func (s StringSlice) ToString(v []string) string {
	if len(v) == 0 {
		return ""
	}

	return fmt.Sprintf("%q", strings.Join(v, s.delimiter))
}

// Set implements the flag.Value interface.
func (s *StringSlice) Set(v string) error {
	if v == "" {
		*s.destination = []string{}

		return nil
	}

	out := strings.Split(v, s.delimiter)

	//nolint:mnd
	for i := len(out) - 2; i >= 0; i-- {
		if strings.HasSuffix(out[i], s.escapeString) {
			out[i] = out[i][:len(out[i])-len(s.escapeString)] + s.delimiter + out[i+1]
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

	return strings.Join(*s.destination, s.delimiter)
}
