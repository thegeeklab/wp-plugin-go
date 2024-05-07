package types

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCmdRun(t *testing.T) {
	tests := []struct {
		name       string
		cmd        *Cmd
		wantErr    bool
		wantStdout string
		wantStderr string
		wantTrace  string
	}{
		{
			name: "trace enabled",
			cmd: &Cmd{
				Trace: boolPtr(true),
				Cmd: &exec.Cmd{
					Path: "/usr/bin/echo",
					Args: []string{"echo", "hello"},
				},
			},
			wantTrace:  "+ echo hello\n",
			wantStdout: "hello\n",
		},
		{
			name: "private output",
			cmd: &Cmd{
				Private: true,
				Cmd: &exec.Cmd{
					Path: "/usr/bin/echo",
					Args: []string{"echo", "hello"},
				},
			},
			wantTrace: "+ echo hello\n",
		},
		{
			name: "custom env",
			cmd: &Cmd{
				Cmd: &exec.Cmd{
					Path: "/bin/sh",
					Args: []string{"sh", "-c", "echo $TEST"},
					Env:  []string{"TEST=1"},
				},
			},
			wantTrace:  "+ sh -c echo $TEST\n",
			wantStdout: "1\n",
		},
		{
			name: "custom stdout",
			cmd: &Cmd{
				Cmd: &exec.Cmd{
					Path:   "/bin/sh",
					Args:   []string{"sh", "-c", "echo hello"},
					Stdout: new(bytes.Buffer),
				},
			},
			wantTrace:  "+ sh -c echo hello\n",
			wantStdout: "hello\n",
		},
		{
			name: "custom stderr",
			cmd: &Cmd{
				Cmd: &exec.Cmd{
					Path:   "/bin/sh",
					Args:   []string{"sh", "-c", "echo error >&2"},
					Stderr: new(bytes.Buffer),
				},
			},
			wantTrace:  "+ sh -c echo error >&2\n",
			wantStderr: "error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			traceBuf := new(bytes.Buffer)
			stdoutBuf := new(bytes.Buffer)
			stderrBuf := new(bytes.Buffer)
			tt.cmd.TraceWriter = traceBuf
			tt.cmd.Stdout = stdoutBuf
			tt.cmd.Stderr = stderrBuf

			err := tt.cmd.Run()
			if tt.wantErr {
				assert.Error(t, err)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantTrace, traceBuf.String())
			assert.Equal(t, tt.wantStdout, stdoutBuf.String())
			assert.Equal(t, tt.wantStderr, stderrBuf.String())
		})
	}
}

func TestCmdSetTrace(t *testing.T) {
	tests := []struct {
		name     string
		cmd      *Cmd
		trace    bool
		expected *bool
	}{
		{
			name:     "set trace to true",
			cmd:      &Cmd{},
			trace:    true,
			expected: boolPtr(true),
		},
		{
			name:     "set trace to false",
			cmd:      &Cmd{},
			trace:    false,
			expected: boolPtr(false),
		},
		{
			name:     "overwrite existing trace value",
			cmd:      &Cmd{Trace: boolPtr(true)},
			trace:    false,
			expected: boolPtr(false),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cmd.SetTrace(tt.trace)
			assert.Equal(t, tt.expected, tt.cmd.Trace)
		})
	}
}

func boolPtr(b bool) *bool {
	return &b
}
