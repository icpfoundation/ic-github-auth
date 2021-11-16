package main

import (
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/lyswifter/ic-auth/db"
	"github.com/lyswifter/ic-auth/server"
	"github.com/lyswifter/ic-auth/util"
	"github.com/mitchellh/go-homedir"
	"github.com/unrolled/secure"
)

func main() {
	util.Infof("This is the internet computer github app authorization server")

	db.DataStores()

	setupAuthServer()
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next()
	}
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

	// r.Run("0.0.0.0:9091") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	certpath, err := homedir.Expand("~/.cert")
	if err != nil {
		return
	}

	r.RunTLS("0.0.0.0:9091", path.Join(certpath, "5537464__skyipfs.com.pem"), path.Join(certpath, "5537464__skyipfs.com.key"))
}
