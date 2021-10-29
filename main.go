package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

const repoPath = "~/.icauth"

const client_id = "Iv1.018aba55453994ac"
const client_secret = "e6a5b65152a4dca9754fa2e13df80f3c087019e7"

var accessTokenUrl = "https://github.com/login/oauth/access_token"
var redirect_uri = "http://54.244.200.160/:9091/public/auth/"
var state = "xxxxxx"

func main() {
	Infof("This is the internet computer github authorize server")

	DataStores()

	setupAuthServer()
}

func setupAuthServer() {
	r := gin.Default()

	handleAccessTokenRedirectAPI(r)

	handleGithubAuthorizeAPI(r)

	r.Run("0.0.0.0:9091") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func handleAccessTokenRedirectAPI(r *gin.Engine) {
	r.GET("/public/token", func(c *gin.Context) {
		Infof("get access token redirect url: %s", c.Request.URL.String())

		r.LoadHTMLFiles("./index.html")
		c.HTML(200, "./index.html", nil)
	})
}

func handleGithubAuthorizeAPI(r *gin.Engine) {
	r.GET("/public/auth", func(c *gin.Context) {
		code := c.Query("code")

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

		r.LoadHTMLFiles("./index.html")
		c.HTML(200, "./index.html", nil)
	})
}

func getAccessToken(code string, redirect_uri string, state string) error {
	url := fmt.Sprintf("%s?client_id=%s&client_secret=%s&code=%s&redirect_uri=%s&state=%s", accessTokenUrl, client_id, client_secret, code, redirect_uri, state)
	method := "POST"

	Infof("access_token url: %s", url)

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

	err = SaveAccessToken(context.TODO(), state, AuthResponse{
		AccessToken:          "",
		ExpiresIn:            121,
		RefreshTOken:         "xxx",
		RefreshTokenExpireIn: 333,
		Scope:                "xxx",
		TokenType:            "xxx",
	})
	if err != nil {
		Errorf("save access token err %s", err.Error())
		return err
	}

	return nil
}
