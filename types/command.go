package types

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/sys/execabs"
)

type Cmd struct {
	*execabs.Cmd
	Private bool
	Trace   *bool
}

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

	if *c.Trace {
		fmt.Fprintf(os.Stdout, "+ %s\n", strings.Join(c.Args, " "))
	}

	if err := c.Start(); err != nil {
		return err
	}

	return c.Wait()
}

func (c *Cmd) SetTrace(trace bool) {
	c.Trace = &trace
}
