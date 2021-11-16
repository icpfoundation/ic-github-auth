package server

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/lyswifter/ic-auth/util"
)

func HandleAccessTokenRedirectAPI(r *gin.Engine) {
	r.GET("/public/token", func(c *gin.Context) {
		util.Infof("get access token url: %s", c.Request.URL.String())
		state := c.Query("state")

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, TRACE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type,Token,Accept, Connection, User-Agent, Cookie")
		c.Header("Access-Control-Max-Age", "3628800")

		if state == "" {
			c.JSON(502, gin.H{
				"status":  "Err",
				"message": "state must not nil",
			})
			return
		}

		code, err := Authdb.ReadInstallCode(context.TODO(), state)
		if err != nil {
			c.JSON(502, gin.H{
				"status":  "Err",
				"message": "could not get installation code",
			})
			return
		}

		ret, err := Authdb.ReadAccessToken(context.TODO(), state)
		if err != nil {
			c.JSON(502, gin.H{
				"status":  "Err",
				"message": "could not get access code",
			})
			return
		}

		if ret == nil {
			c.JSON(502, gin.H{
				"status":  "Err",
				"message": "access token is nil",
			})
			return
		}

		util.Infof("read installation code: %s token: %s", string(code), string(ret))

		c.JSON(200, gin.H{
			"statue":  "Ok",
			"message": "success",
			"token":   string(ret),
			"code":    string(code),
		})
	})
}
