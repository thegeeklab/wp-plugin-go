package exec

import (
	"bytes"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCmdRun(t *testing.T) {
	echoPath, err := exec.LookPath("echo")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		cmd        *Cmd
		wantErr    error
		wantStdout string
		wantStderr string
		wantTrace  string
	}{
		{
			name: "trace enabled",
			cmd: &Cmd{
				Trace: true,
				Cmd: &exec.Cmd{
					Path: echoPath,
					Args: []string{"echo", "hello"},
				},
			},
			wantTrace:  "▶  echo hello\n",
			wantStdout: "hello\n",
		},
		{
			name: "trace legacy",
			cmd: &Cmd{
				Trace:       true,
				TraceLegacy: true,
				Cmd: &exec.Cmd{
					Path: echoPath,
					Args: []string{"echo", "hello"},
				},
			},
			wantTrace:  "+ echo hello\n",
			wantStdout: "hello\n",
		},
		{
			name: "trace disabled",
			cmd: &Cmd{
				Trace: false,
				Cmd: &exec.Cmd{
					Path: echoPath,
					Args: []string{"echo", "hello"},
				},
			},
			wantStdout: "hello\n",
		},
		{
			name: "custom env",
			cmd: &Cmd{
				Trace: true,
				Cmd: &exec.Cmd{
					Path: "/bin/sh",
					Args: []string{"sh", "-c", "echo $TEST"},
					Env:  []string{"TEST=1"},
				},
			},
			wantTrace:  "▶  sh -c echo $TEST\n",
			wantStdout: "1\n",
		},
		{
			name: "custom stderr",
			cmd: &Cmd{
				Trace: true,
				Cmd: &exec.Cmd{
					Path:   "/bin/sh",
					Args:   []string{"sh", "-c", "echo error >&2"},
					Stderr: new(bytes.Buffer),
				},
			},
			wantTrace:  "▶  sh -c echo error >&2\n",
			wantStderr: "error\n",
		},
		{
			name: "error",
			cmd: &Cmd{
				Trace: true,
				Cmd: &exec.Cmd{
					Path: "/invalid/path",
				},
			},
			wantErr: assert.AnError,
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
			if tt.wantErr != nil {
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

func TestCommand(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		args     []string
		wantArgs []string
	}{
		{
			name:     "command without arguments",
			command:  "echo",
			args:     []string{},
			wantArgs: []string{"echo"},
		},
		{
			name:     "command with arguments",
			command:  "echo",
			args:     []string{"hello", "world"},
			wantArgs: []string{"echo", "hello", "world"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Command(tt.command, tt.args...) //nolint:gosec // command is only constructed, not executed

			assert.True(t, got.Trace)
			assert.Equal(t, os.Stdout, got.TraceWriter)
			assert.Equal(t, os.Environ(), got.Env)
			assert.NotEmpty(t, got.Path)
			assert.Equal(t, tt.wantArgs, got.Args)
		})
	}
}
