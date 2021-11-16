package cmd

import "github.com/urfave/cli"

var DeployCmd = cli.Command{
	Name:        "deploy",
	Description: "Deploy app to Internet Computer",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "path",
		},
	},
	Action: func(cctx *cli.Context) error {

		return nil
	},
}
