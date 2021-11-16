package main

import (
	"fmt"
	"os"

	"github.com/lyswifter/ic-auth/cmd"
	"github.com/lyswifter/ic-auth/util"
	"github.com/urfave/cli"
)

func main() {
	util.Infof("This is the internet computer github app authorization server")

	local := []cli.Command{
		cmd.InitCmd,
		cmd.DaemonCmd,
		cmd.DeployCmd,
	}

	app := &cli.App{
		Name:    "ic-auth",
		Usage:   "Used for authorize user to access Internet Computer",
		Version: "0.0.1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "icAuth-dir",
				Value: "~/.icauth",
			},
		},

		Commands: local,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
