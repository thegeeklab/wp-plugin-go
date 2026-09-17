package flags

import "github.com/urfave/cli/v3"

func flags() []cli.Flag {
	return []cli.Flag{
		// Dummy flag long description spanning
		// two source lines in the same paragraph.
		//
		// Second paragraph of the dummy flag long description.
		&cli.StringFlag{
			Name:     "dummy-flag",
			Usage:    "Dummy flag desc.",
			Sources:  cli.EnvVars("PLUGIN_DUMMY_FLAG"),
			Required: true,
		},

		&cli.IntFlag{
			Name:     "dummy-flag-int",
			Usage:    "dummy int flag desc",
			Sources:  cli.EnvVars("PLUGIN_DUMMY_FLAG_INT"),
			Required: true,
		},

		// Long description for the slice flag with
		// multiple paragraphs.
		//
		// Second paragraph for slice flag.
		&cli.StringSliceFlag{
			Name:    "slice.flag",
			Usage:   "slice flag",
			Sources: cli.EnvVars("PLUGIN_SLICE_FLAG"),
		},

		&cli.StringFlag{
			Name:    "simpe-flag",
			Value:   "simple",
			Sources: cli.EnvVars("PLUGIN_X_SIMPLE_FLAG"),
		},

		&cli.StringFlag{
			Name:    "other.flag",
			Usage:   "other flag with desc",
			Sources: cli.EnvVars("PLUGIN_Z_OTHER_FLAG"),
		},

		&cli.StringFlag{
			Name:    "hidden.flag",
			Usage:   "hidden flag",
			Sources: cli.EnvVars("HIDDEN_FLAG", "PLUGIN_HIDDEN_FLAG"),
		},
	}
}
