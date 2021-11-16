package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lyswifter/ic-auth/db"
	"github.com/lyswifter/ic-auth/server"
	"github.com/lyswifter/ic-auth/util"
	"github.com/mitchellh/go-homedir"
	"github.com/unrolled/secure"
	"github.com/urfave/cli"
)

var DaemonCmd = cli.Command{
	Name:        "daemon",
	Description: "Start ic auth daemon",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "repo",
			Value: "~/.icauth",
			Usage: "Specify the location of database",
		},
	},
	Action: func(cctx *cli.Context) error {
		repodir, err := homedir.Expand(cctx.String("repo"))
		if err != nil {
			return err
		}

		authdb, err := db.SetupAuth(repodir)
		if err != nil {
			return err
		}
		server.Authdb = authdb

		fmt.Printf("server.Authdb: %+v\n", server.Authdb)

		go setupAuthServer()

		ticker := time.NewTicker(30 * time.Second)
		ctx := context.Background()
		sigChan := make(chan os.Signal, 2)

		util.Infof("Internet Computer auth is running...")

		for {
			select {
			case <-ticker.C:
				util.Infof("I am running")
			case <-ctx.Done():
				util.Infof("Shutting down..")
				util.Infof("Graceful shutdown successful")
				signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
				return nil
			}
		}
	},
}

func setupAuthServer() {
	secureMiddleware := secure.New(secure.Options{
		FrameDeny: true,
	})

	secureFunc := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			err := secureMiddleware.Process(c.Writer, c.Request)

			// If there was an error, do not continue.
			if err != nil {
				c.Abort()
				return
			}

			// Avoid header rewrite if response is a redirection.
			if status := c.Writer.Status(); status > 300 && status < 399 {
				c.Abort()
			}
		}
	}()

	r := gin.Default()
	r.Use(secureFunc)

	server.HandleAccessTokenRedirectAPI(r)
	server.HandleGithubAuthorizeAPI(r)
	server.HandleTiggerBuildAPI(r)
	server.HandleDeployLogAPI(r)
	server.HandleCanisterListAPI(r)
	server.HandleCanisterInfoAPI(r)

	r.Run("0.0.0.0:9091") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	// certpath, err := homedir.Expand("~/.cert")
	// if err != nil {
	// 	return
	// }

	// r.RunTLS("0.0.0.0:9091", path.Join(certpath, "5537464__skyipfs.com.pem"), path.Join(certpath, "5537464__skyipfs.com.key"))
}
