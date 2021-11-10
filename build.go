package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/go-homedir"
)

var addr = "http://chaincloud.skyipfs.com:9091"

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

		logpath := path.Join(repo, "repository", reponame, "logs")
		err = mkDir(logpath)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		// 3. clone target repo and branch source code to target directory
		clonecmd := exec.Command("git", "clone", "-b", branch, repourl, targetpath)
		clonecmd.Stderr = os.Stderr
		clonecmd.Stdout = os.Stdout
		err = clonecmd.Run()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		var connectionId = time.Now().Unix()

		switch framework {
		case "dfx":
			deploycmd := exec.Command("dfx", "deploy", "--network", "ic")
			deploycmd.Dir = targetpath

			stderr, err := deploycmd.StderrPipe()
			if err != nil {
				return
			}

			stdout, err := deploycmd.StdoutPipe()
			if err != nil {
				return
			}

			deploycmd.Start()

			errReader := bufio.NewReader(stderr)
			outReader := bufio.NewReader(stdout)

			fmt.Printf("conn: %+v", c)

			f, err := os.Create(path.Join(logpath, fmt.Sprintf("%d", timing)))
			if err != nil {
				return
			}

			defer f.Close()

			for {
				line, err := errReader.ReadString('\n')
				if err == io.EOF {
					break
				}

				if err != nil {
					break
				}

				// write local
				_, err = f.WriteString(line)
				if err != nil {
					break
				}
			}

			for {
				line, err := outReader.ReadString('\n')
				if err == io.EOF {
					break
				}

				if err != nil {
					break
				}

				// write local
				_, err = f.WriteString(line)
				if err != nil {
					break
				}
			}

			deploycmd.Wait()
		default:
		}

		c.JSON(http.StatusOK, gin.H{
			"statue":       "Ok",
			"message":      "tigger build ok",
			"connectionid": connectionId,
		})
	})
}

func startLocalNetworkWithDfx(path string) ([]byte, error) {
	// 3. start local network
	startcmd := exec.Command("dfx", "start", "--background")
	startcmd.Dir = path

	var b bytes.Buffer
	startcmd.Stdout = &b
	startcmd.Stderr = &b
	err := startcmd.Run()
	if err != nil {
		fmt.Printf("dfx(%s) err: %s ret: %s\n", path, err.Error(), b.String())
		return b.Bytes(), err
	}
	return b.Bytes(), nil
}

func deployWithHugo(path string) {

}

///////////////////////////////
