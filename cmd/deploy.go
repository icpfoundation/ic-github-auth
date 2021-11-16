package cmd

import "github.com/urfave/cli"

var DeployCmd = cli.Command{
	Name:        "deploy",
	Description: "Deploy app to Internet Computer",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "type",
			Value: "dfx",
			Usage: "Specify frontend framework to use [reactjs, nuxtjs, nextjs, vuejs, hugo]",
		},
	},
	Action: func(cctx *cli.Context) error {
		//

		return nil
	},
}
