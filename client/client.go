package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}

	token, err := getToken()
	if err != nil {
		fmt.Printf("Error when getting token: %s", err.Error())
	}

	getResource("", "/api/public")
	getResource(token, "/api/private")
	getResource(token, "/api/private-scoped")
}

func getToken() (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}

	url := "https://" + os.Getenv("AUTH0_DOMAIN") + "/oauth/token"
	payload := strings.NewReader("{\"client_id\":\"" + os.Getenv("CLIENT_ID") + "\",\"client_secret\":\"" + os.Getenv("CLIENT_SECRET") + "\",\"audience\":\"" + os.Getenv("AUTH0_AUDIENCE") + "\",\"grant_type\":\"client_credentials\"}")

	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	var tr = tokenResponse{}
	err = json.NewDecoder(res.Body).Decode(&tr)

	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s %s", tr.TokenType, tr.AccessToken), nil
}

func getResource(token, relEnd string) {

	url := "http://localhost:3010" + relEnd
	req, _ := http.NewRequest("GET", url, nil)

	if token != "" {
		req.Header.Add("authorization", token)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}
