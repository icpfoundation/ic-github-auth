package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lyswifter/ic-auth/params"
	"github.com/lyswifter/ic-auth/util"
)

func HandleGithubAuthorizeAPI(r *gin.Engine) {
	r.GET("/public/auth", func(c *gin.Context) {
		code := c.Query("code")
		state := c.Query("state")
		installationId := c.Query("installation_id")

		util.Infof("get authorize redirect(%s): %s", c.Request.URL.String(), code)

		if code == "" {
			c.JSON(500, gin.H{
				"status":  "Err",
				"message": "could not get auth code",
			})
			return
		}

		err := Authdb.SaveInstallCode(context.TODO(), state, installationId)
		if err != nil {
			return
		}

		err = getAccessToken(code, params.Redirect_uri, state)
		if err != nil {
			util.Errorf("get access token: %s", err.Error())
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
	url := fmt.Sprintf("%s?client_id=%s&client_secret=%s&code=%s&redirect_uri=%s&state=%s", params.AccessTokenUrl, params.Client_id, params.Client_secret, code, redirect_uri, state)
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		util.Errorf("new request err: %s", err.Error())
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		util.Errorf("do request err: %s", err.Error())
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		util.Errorf("read response err: %s", err.Error())
		return err
	}
	util.Infof("request access token response: %s", string(body))

	//access_token=ghu_xaM3AEToCoATeS4S0lT3hyxZfXu9Jg4Wwomz&expires_in=28800&
	//refresh_token=ghr_MYVUlf0gKOpucwkPUWoxYpjSzRu2Rpn5Gy9pzIPD0DA8LgJNYYkGLNG6OusAGzrsa5GNaT0k8Z9d&refresh_token_expires_in=15724800&scope=&token_type=bearer

	if state == "" {
		state = "xxxxxxx"
	}

	err = Authdb.SaveAccessToken(context.TODO(), state, body)
	if err != nil {
		util.Errorf("save access token err %s", err.Error())
		return err
	}

	return nil
}
