package main

import (
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/go-homedir"
	"github.com/unrolled/secure"
)

const repoPath = "~/.icauth"

const client_id = "Iv1.018aba55453994ac"
const client_secret = "e6a5b65152a4dca9754fa2e13df80f3c087019e7"

var accessTokenUrl = "https://github.com/login/oauth/access_token"
var redirect_uri = "http://54.244.200.160:9091/public/auth/"

func main() {
	Infof("This is the internet computer github app authorization server")

	DataStores()

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

	handleAccessTokenRedirectAPI(r)
	handleGithubAuthorizeAPI(r)
	handleTiggerBuildAPI(r)

	// r.Run("0.0.0.0:9091") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	certpath, err := homedir.Expand("~/.cert")
	if err != nil {
		return
	}

	r.RunTLS("0.0.0.0:9091", path.Join(certpath, "5537464__skyipfs.com.pem"), path.Join(certpath, "5537464__skyipfs.com.key"))
}
