package exec

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
				Trace: true,
				Cmd: &exec.Cmd{
					Path: "/usr/bin/echo",
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
					Path: "/usr/bin/echo",
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
			wantTrace:  "+ sh -c echo $TEST\n",
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
			wantTrace:  "+ sh -c echo error >&2\n",
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
			wantErr: true,
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
