package types

import "golang.org/x/sys/execabs"

type Cmd struct {
	*execabs.Cmd
	Private bool
}
