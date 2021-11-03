package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/go-homedir"
)

func handleTiggerBuildAPI(r *gin.Engine) {
	r.GET("public/build", func(c *gin.Context) {
		// 1. get target repo url and branch
		// 2. clone source code to specify directory
		// 3. if specify website generator, than using specify build command and move ouput files to /dist directory
		// 4. if no canister on the mainnet, than generate canister firstly and topup some cycle into it
		// 5. run dfx build and dfx install / dfx deploy to deploy canisters to mainnet
		// 6. get deploy process and status output file
		// 7. get deploy canister id

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, TRACE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type,Token,Accept, Connection, User-Agent, Cookie")
		c.Header("Access-Control-Max-Age", "3628800")

		// 1. parse params
		framework := c.Query("framework")
		reponame := c.Query("reponame")
		repourl := c.Query("repourl")
		branch := c.Query("branch")

		Infof("Tigger build from client for %s %s %s %s", repourl, branch, reponame, framework)

		// 2. mkdir
		timing := time.Now().Unix()
		repo, err := homedir.Expand(repoPath)
		if err != nil {
			return
		}
		targetpath := path.Join(repo, "repository", reponame, fmt.Sprintf("%d", timing))
		err = mkDir(targetpath)
		if err != nil {
			return
		}

		// 3. clone target repo and branch source code to target directory
		clonecmd := exec.Command("git", "clone", "-b", branch, repourl, targetpath)
		fmt.Printf("clonecmd: %v\n", clonecmd)
		clonecmd.Stdout = os.Stdout
		clonecmd.Stderr = os.Stderr
		err = clonecmd.Run()
		if err != nil {
			fmt.Printf("err: %s\n", err.Error())
			return
		}

		// 4. if using default dfx to create a canister
		deploycmd := exec.Command("dfx", "deploy", "--network", "ic")
		deploycmd.Dir = targetpath
		result, err := deploycmd.CombinedOutput()
		if err != nil {
			fmt.Printf("err: %s\n", err.Error())
			return
		}

		Infof("deploy log: %s", string(result))

		// n. recall
		c.String(http.StatusOK, "tigger build successfully")
	})
}
