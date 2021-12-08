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

func HandleRefreshTokenAPI(r *gin.Engine) {
	r.GET("/public/refresh", func(c *gin.Context) {
		refreshToken := c.Query("refresh_token")
		clientId := c.Query("client_id")
		clientSecret := c.Query("client_secret")

		grant_type := "refresh_token"

		ret, err := refreshAccessToken(refreshToken, grant_type, clientId, clientSecret)
		if err != nil {
			c.JSON(500, gin.H{
				"status":  "Err",
				"message": "could not refresh token",
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

func refreshAccessToken(refresh_token string, grant_type string, client_id string, client_secret string) (string, error) {
	url := fmt.Sprintf("%s?refresh_token=%s&grant_type=%s&client_id=%s&client_secret=%s", params.AccessTokenUrl, refresh_token, grant_type, client_id, client_secret)
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		util.Errorf("new request err: %s", err.Error())
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		util.Errorf("do request err: %s", err.Error())
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		util.Errorf("read response err: %s", err.Error())
		return "", err
	}
	util.Infof("refresh access token response: %s", string(body))

	// err = Authdb.SaveAccessToken(context.TODO(), state, body)
	// if err != nil {
	// 	util.Errorf("save access token err %s", err.Error())
	// 	return err
	// }

	return string(body), nil
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
