package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

func setupAuthServer() {
	r := gin.Default()

	r.Use(cors.Default())

	handleAccessTokenRedirectAPI(r)
	handleGithubAuthorizeAPI(r)

	r.Run("0.0.0.0:9091") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func handleAccessTokenRedirectAPI(r *gin.Engine) {
	r.GET("/public/token", func(c *gin.Context) {
		Infof("get access token url: %s", c.Request.URL.String())
		state := c.Query("state")

		if state == "" {
			c.JSON(502, gin.H{
				"status":  "Err",
				"message": "state must not nil",
			})
		}

		ret, err := ReadAccessToken(context.TODO(), state)
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

		c.JSON(200, gin.H{
			"statue":  "Ok",
			"message": "success",
			"token":   string(ret),
		})
	})
}

func handleGithubAuthorizeAPI(r *gin.Engine) {
	r.GET("/public/auth", func(c *gin.Context) {
		code := c.Query("code")
		state := c.Query("state")

		Infof("get authorize redirect(%s): %s", c.Request.URL.String(), code)

		if code == "" {
			c.JSON(500, gin.H{
				"status":  "Err",
				"message": "could not get auth code",
			})
			return
		}

		err := getAccessToken(code, redirect_uri, state)
		if err != nil {
			Errorf("get access token: %s", err.Error())
			return
		}

		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, `<!DOCTYPE html>
		<html lang="en">
		
		<head>
			<meta charset="UTF-8">
			<meta http-equiv="X-UA-Compatible" content="IE=edge">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Chain-cloud</title>
		
			<style type="text/css">
				.bg {
					width: 20%;
					margin: 300px auto;
					text-align: center;
				}
			</style>

			<script>
				window.setTimeout(() => {
					window.close();
				}, 2000)
    		</script>
		</head>
		
		<body>
			<div class="bg">
				<img src="https://storageapi.fleek.co/lyswifter-team-bucket/chain-cloud/nav_logo@2x.png" alt="logo">
				<h3 class="authed">Authorized</h3>
			</div>
		</body>
		
		</html>`)
	})
}

func getAccessToken(code string, redirect_uri string, state string) error {
	url := fmt.Sprintf("%s?client_id=%s&client_secret=%s&code=%s&redirect_uri=%s&state=%s", accessTokenUrl, client_id, client_secret, code, redirect_uri, state)
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		Errorf("new request err: %s", err.Error())
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		Errorf("do request err: %s", err.Error())
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		Errorf("read response err: %s", err.Error())
		return err
	}
	Infof("request access token response: %s", string(body))

	//access_token=ghu_xaM3AEToCoATeS4S0lT3hyxZfXu9Jg4Wwomz&expires_in=28800&
	//refresh_token=ghr_MYVUlf0gKOpucwkPUWoxYpjSzRu2Rpn5Gy9pzIPD0DA8LgJNYYkGLNG6OusAGzrsa5GNaT0k8Z9d&refresh_token_expires_in=15724800&scope=&token_type=bearer

	if state == "" {
		state = "xxxxxxx"
	}

	err = SaveAccessToken(context.TODO(), state, body)
	if err != nil {
		Errorf("save access token err %s", err.Error())
		return err
	}

	return nil
}
