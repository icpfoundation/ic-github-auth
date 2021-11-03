package main

import (
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func handleTiggerBuildAPI(r *gin.Engine) {
	r.GET("public/build", func(c *gin.Context) {
		//1. get target repo url and branch
		//2. clone source code to specify directory
		//3. if specify website generator, than using specify build command and move ouput files to /dist directory
		//4. if no canister on the mainnet, than generate canister firstly and topup some cycle into it
		//5. run dfx build and dfx install / dfx deploy to deploy canisters to mainnet
		//6. get deploy process and status output file
		//7. get deploy canister id

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, TRACE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type,Token,Accept, Connection, User-Agent, Cookie")
		c.Header("Access-Control-Max-Age", "3628800")

		repo := c.Query("repo")
		branch := c.Query("branch")
		Infof("Tigger build from client for %s and %s", repo, branch)

		clonecmd := exec.Command("git", "clone", repo)
		clonecmd.Stdout = os.Stdout
		clonecmd.Stderr = os.Stderr
		err := clonecmd.Run()
		if err != nil {
			Errorf("err: %s", err.Error())
			return
		}

		c.String(http.StatusOK, "tigger build")
	})
}
