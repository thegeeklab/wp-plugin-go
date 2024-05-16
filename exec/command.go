package exec

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sys/execabs"
)

// Cmd represents a command to be executed, with options to control its behavior.
// The Cmd struct embeds the standard library's exec.Cmd, adding additional fields
// to control the command's output and tracing.
type Cmd struct {
	*exec.Cmd
	Trace       bool      // Print composed command before execution.
	TraceWriter io.Writer // Where to write the trace output.
}

// Run runs the command and waits for it to complete.
// If there is an error starting the command, it is returned.
// Otherwise, the command is waited for and its exit status is returned.
func (c *Cmd) Run() error {
	if c.Trace {
		fmt.Fprintf(c.TraceWriter, "+ %s\n", strings.Join(c.Args, " "))
	}

	if err := c.Start(); err != nil {
		return err
	}

	return c.Wait()
}

// Command creates a new Cmd with the given name and arguments. The Cmd is configured
// to use the current environment, and to write its stdout and stderr to the
// process's stdout and stderr.
func Command(name string, arg ...string) *Cmd {
	cmd := &Cmd{
		Cmd:         execabs.Command(name, arg...),
		Trace:       true,
		TraceWriter: os.Stdout,
	}

	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}
