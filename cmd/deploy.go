package cmd

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/lyswifter/ic-auth/db"
	"github.com/lyswifter/ic-auth/deploy"
	"github.com/lyswifter/ic-auth/server"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

var DeployCmd = cli.Command{
	Name:        "deploy",
	Description: "Deploy app to Internet Computer",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "db",
			Value: "~/.icauth",
			Usage: "Specify the location of database",
		},
		&cli.StringFlag{
			Name:  "framework",
			Value: "dfx",
			Usage: "Specify frontend framework to use [reactjs, nuxtjs, nextjs, vuejs, hugo]",
		},
		&cli.StringFlag{
			Name:  "repo",
			Usage: "Specify the repo name to deploy",
		},
		&cli.StringFlag{
			Name:  "target",
			Usage: "Specify target path to minupate",
		},
		&cli.BoolFlag{
			Name:  "islocal",
			Usage: "Speicfy whether is local network",
		},
		&cli.StringFlag{
			Name:  "output",
			Usage: "Specify output log file location",
		},
		&cli.StringFlag{
			Name:  "cname",
			Usage: "Specify canister name if there is no one exist",
		},
		&cli.StringFlag{
			Name:  "outsource",
			Usage: "Specify build output source file location",
		},
	},
	Action: func(cctx *cli.Context) error {
		framework := cctx.String("framework")
		repo := cctx.String("repo")
		target := cctx.String("target")
		islocal := cctx.Bool("islocal")
		fmt.Printf("Will deploy repo: %s framework: %s target: %s islocal: %t\n", repo, target, framework, islocal)

		canistername := cctx.String("cname")
		outsource := cctx.String("outsource")
		fmt.Printf("canistername: %s outsource: %s\n", canistername, outsource)

		logpath, err := homedir.Expand("~/log")
		if err != nil {
			return err
		}

		timing := time.Now().Unix()
		f, err := os.Create(path.Join(logpath, fmt.Sprintf("%s-%d.out", repo, timing)))
		if err != nil {
			return err
		}

		authdb, err := db.SetupAuth(cctx.String("db"))
		if err != nil {
			return err
		}

		server.Authdb = authdb

		switch framework {
		case "dfx":
			cinfos, err := deploy.DeployWithDfx(target, f, repo, islocal, framework, "")
			if err != nil {
				return err
			}

			for _, v := range cinfos {
				err := authdb.SaveCanisterInfo(context.TODO(), v)
				if err != nil {
					return err
				}
			}
		case "reactjs":
			cinfos, err := deploy.DeployWithReactjs(target, f, canistername, outsource, repo, islocal, framework)
			if err != nil {
				return err
			}

			for _, v := range cinfos {
				err := authdb.SaveCanisterInfo(context.TODO(), v)
				if err != nil {
					return err
				}
			}
		case "nuxtjs":
			cinfos, err := deploy.DeployWithNuxt(target, f, canistername, outsource, repo, islocal, framework)
			if err != nil {
				return err
			}

			for _, v := range cinfos {
				err := authdb.SaveCanisterInfo(context.TODO(), v)
				if err != nil {
					return err
				}
			}
		default:
		}

		setupAuthServer()

		return nil
	},
}
