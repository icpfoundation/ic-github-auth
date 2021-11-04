package main

import (
	"bytes"
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
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		targetpath := path.Join(repo, "repository", reponame, fmt.Sprintf("%d", timing))
		err = mkDir(targetpath)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		// 3. clone target repo and branch source code to target directory
		clonecmd := exec.Command("git", "clone", "-b", branch, repourl, targetpath)
		fmt.Printf("clonecmd: %v\n", clonecmd)
		clonecmd.Stderr = os.Stderr
		clonecmd.Stdout = os.Stdout
		err = clonecmd.Run()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		var retLog string
		switch framework {
		case "dfx":
			ret, err := deployWithDfx(targetpath)
			if err != nil {
				retbyte, err := buildOutLogs(string(ret))
				if err != nil {
					return
				}

				c.String(http.StatusAccepted, string(retbyte))
				return
			}
			retLog = string(ret)
		default:
		}

		// Infof("retLog: %d", len(strings.Split(retLog, "\n")))

		retbyte, err := buildOutLogs(string(retLog))
		if err != nil {
			return
		}

		// n. recall
		c.String(http.StatusOK, string(retbyte))
	})
}

func deployWithDfx(path string) ([]byte, error) {
	// 4. if using default dfx to create a canister
	deploycmd := exec.Command("dfx", "deploy", "--network", "ic")
	deploycmd.Dir = path

	var b bytes.Buffer
	deploycmd.Stdout = &b
	deploycmd.Stderr = &b
	err := deploycmd.Run()
	if err != nil {
		fmt.Printf("dfx(%s) err: %s ret: %s\n", path, err.Error(), b.String())
		return b.Bytes(), err
	}

	return b.Bytes(), nil
}

func deployWithHugo(path string) {

}
