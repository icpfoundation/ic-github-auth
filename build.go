package main

import (
	"bufio"
	"bytes"
	"errors"
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

func handleDeployLogAPI(r *gin.Engine) {
	r.GET("public/logs", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, TRACE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type,Token,Accept, Connection, User-Agent, Cookie")
		c.Header("Access-Control-Max-Age", "3628800")

		startline := c.Query("startline")
		endline := c.Query("endline")

		filename := c.Query("file")
		if filename == "" {
			c.String(http.StatusInternalServerError, errors.New("filename must provide").Error())
		}

		reponame := c.Query("reponame")

		///
		repo, err := homedir.Expand(repoPath)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		filepath := path.Join(repo, "repository", reponame, "logs", filename)
		fmt.Printf("filepath: %s\n", filepath)

		line := fmt.Sprintf("%s,%sp;", startline, endline)
		tailCmd := exec.Command("sed", "-n", line, filepath)

		var b bytes.Buffer
		tailCmd.Stderr = &b
		tailCmd.Stdout = &b

		tailCmd.Run()
		if err != nil {
			fmt.Printf("tail err: %s", err.Error())
			return
		}

		c.String(http.StatusOK, b.String())
	})
}

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

		switch framework {
		case "dfx":

			go func() {
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

				f, err := os.Create(path.Join(logpath, fmt.Sprintf("%d", timing)))
				if err != nil {
					return
				}

				defer f.Close()

				// defer func() {
				// 	err := cleanCacheCmd(targetpath)
				// 	if err != nil {
				// 		return
				// 	}

				// 	fmt.Printf("clean cache at: %s ok", targetpath)
				// }()

				for {
					line, err := errReader.ReadString('\n')
					if err == io.EOF {
						break
					}

					if err != nil {
						break
					}

					// write local
					_, err = f.WriteString(fmt.Sprintf("[%s]	%s", time.Now().Format("2006-01-02 15:04:05.999"), line))
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
					_, err = f.WriteString(fmt.Sprintf("[%s]	%s", time.Now().Format("2006-01-02 15:04:05.999"), line))
					if err != nil {
						break
					}
				}

				deploycmd.Wait()
			}()
		default:
		}

		c.JSON(http.StatusOK, gin.H{
			"statue":       "Ok",
			"message":      "tigger build ok",
			"connectionid": timing,
		})
	})
}

func cleanCacheCmd(path string) error {
	cleanCmd := exec.Command("dfx", "cache", "delete")
	cleanCmd.Dir = path

	cleanCmd.Stderr = os.Stderr
	cleanCmd.Stdout = os.Stdout
	err := cleanCmd.Run()
	if err != nil {
		return err
	}

	return nil
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
