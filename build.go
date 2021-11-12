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

func handleDeployLogAPI(r *gin.Engine) {
	r.GET("public/logs", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, TRACE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type,Token,Accept, Connection, User-Agent, Cookie")
		c.Header("Access-Control-Max-Age", "3628800")

		filename := c.Query("file")
		if filename == "" {
			c.String(http.StatusInternalServerError, "filename must provide")
		}
		reponame := c.Query("reponame")

		repo, err := homedir.Expand(repoPath)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		filepath := path.Join(repo, "logs", reponame, filename)

		ret, err := os.ReadFile(filepath)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.String(http.StatusOK, string(ret))
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
		targetpath := path.Join(repo, "repository", reponame)
		err = mkDir(targetpath)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		logpath := path.Join(repo, "logs", reponame)
		err = mkDir(logpath)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		f, err := os.Create(path.Join(logpath, fmt.Sprintf("%d", timing)))
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		if Exists(path.Join(targetpath, ".git")) {
			pullcmd := exec.Command("git", "pull")
			pullcmd.Dir = targetpath
			pullcmd.Stderr = os.Stderr
			pullcmd.Stdout = os.Stdout
			err = pullcmd.Run()
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
		} else {
			// 3. clone target repo and branch source code to target directory
			clonecmd := exec.Command("git", "clone", "-b", branch, repourl, targetpath)
			clonecmd.Stderr = os.Stderr
			clonecmd.Stdout = os.Stdout
			err = clonecmd.Run()
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
		}

		switch framework {
		case "dfx":
			go func() error {

				defer f.Close()

				err := npmInstall(targetpath, f)
				if err != nil {
					return err
				}

				err = deployWithDfx(targetpath, f)
				if err != nil {
					return err
				}

				defer fmt.Errorf("deploy with dfx: %s", err.Error())

				return nil
			}()

		case "reactjs":
			canistername := c.Query("canistername")
			resource := c.Query("resourcepath")

			if canistername == "" || resource == "" {
				c.String(http.StatusBadRequest, "canister name and resource path must provide")
				return
			}

			go func() error {
				defer f.Close()

				//npm install and npm run build
				err := deployWithReactjs(targetpath, f, canistername, resource)
				if err != nil {
					return err
				}

				defer fmt.Errorf("deploy with reactjs: %s", err.Error())

				return nil
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
