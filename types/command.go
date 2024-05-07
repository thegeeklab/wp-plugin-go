package types

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/sys/execabs"
)

// Cmd represents a command to be executed. It extends the execabs.Cmd struct
// and adds fields for controlling the command's execution, such as whether
// it should be executed in private mode (with output discarded) and whether
// its execution should be traced.
type Cmd struct {
	*execabs.Cmd
	Private     bool
	Trace       *bool
	TraceWriter io.Writer
}

// Run executes the command and waits for it to complete.
// If the Trace field is nil, it is set to true.
// The Env, Stdout, and Stderr fields are set to their default values if they are nil.
// If the Private field is true, the Stdout field is set to io.Discard.
// If the Trace field is true, the command line is printed to Stdout.
func (c *Cmd) Run() error {
	if c.Trace == nil {
		c.SetTrace(true)
	}

	if c.Env == nil {
		c.Env = os.Environ()
	}

	if c.Stdout == nil {
		c.Stdout = os.Stdout
	}

	if c.Stderr == nil {
		c.Stderr = os.Stderr
	}

	if c.Private {
		c.Stdout = io.Discard
	}

	if c.TraceWriter == nil {
		c.TraceWriter = os.Stdout
	}

	if *c.Trace {
		fmt.Fprintf(c.TraceWriter, "+ %s\n", strings.Join(c.Args, " "))
	}

	if err := c.Start(); err != nil {
		return err
	}

	return c.Wait()
}

// SetTrace sets the Trace field of the Cmd to the provided boolean value.
func (c *Cmd) SetTrace(trace bool) {
	c.Trace = &trace
}
