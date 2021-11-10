package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
			http.HandleFunc(fmt.Sprintf("%s/%d", "/public/log", connectionId), func(rw http.ResponseWriter, r *http.Request) {
				var upgrader = websocket.Upgrader{}

				c, err := upgrader.Upgrade(rw, r, nil)
				if err != nil {
					return
				}

				defer c.Close()

				// 4. if using default dfx to create a canister
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

				for {
					line, err := errReader.ReadBytes('\n')
					if err == io.EOF {
						break
					}

					if err != nil {
						break
					}

					err = c.WriteMessage(websocket.TextMessage, line)
					if err != nil {
						log.Println("write stderr:", err)
						break
					}
				}

				for {
					line, err := outReader.ReadBytes('\n')
					if err == io.EOF {
						break
					}

					if err != nil {
						break
					}

					err = c.WriteMessage(websocket.TextMessage, line)
					if err != nil {
						log.Println("write stdout:", err)
						break
					}
				}

				deploycmd.Wait()
			})

			http.ListenAndServe(addr, nil)
		default:
		}

		fmt.Printf("reach here")

		// n. recall
		c.String(http.StatusOK, fmt.Sprintf("tigger build ok with: %d", connectionId))
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
