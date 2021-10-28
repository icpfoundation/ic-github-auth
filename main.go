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
var redirect_uri = ""
var state = ""

func main() {
	Infof("This is the internet computer github authorize server")

	DataStores()

	setupAuthServer()
}

func setupAuthServer() {
	r := gin.Default()

	handleGithubAuthorizeAPI(r)

	r.Run("0.0.0.0:9091") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func handleGithubAuthorizeAPI(r *gin.Engine) {
	r.GET("/public/auth", func(c *gin.Context) {
		code := c.Query("code")

		getAccessToken(code, redirect_uri, state)
	})
}

func getAccessToken(code string, redirect_uri string, state string) {
	url := fmt.Sprintf("%s?client_id=%s&client_secret=%s&code=%s&redirect_uri=%s&state=%s", accessTokenUrl, client_id, client_secret, code, redirect_uri, state)
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

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
		return
	}
}
