package auth

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func Oauth2(w http.ResponseWriter, r *http.Request) {
	clientID := "8134856"
	redirectURI := "https://gamersgazette.herokuapp.com/auth/me"
	scope := []string{"account"}
	state := "12345"
	scopeTemp := strings.Join(scope, "+")
	url := fmt.Sprintf("https://oauth.vk.com/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=%s", clientID, redirectURI, scopeTemp, state)
	fmt.Fprint(w, url)
}

func Me(w http.ResponseWriter, r *http.Request) {
	clientID := "8134856"
	redirectURI := "https://gamersgazette.herokuapp.com/auth/me"
	clientSecret := "7Vw4ALUIHMLPpHTKiRlG"
	state := "12345"
	stateTemp := r.URL.Query().Get("state")
	if stateTemp[len(stateTemp)-1] == '}' {
		stateTemp = stateTemp[:len(stateTemp)-1]
	}
	if stateTemp == "" {
		respErr(w, fmt.Errorf("state query param is not provided"))
		return
	} else if stateTemp != state {
		respErr(w, fmt.Errorf("state query param do not match original one, got=%s", stateTemp))
		return
	}
	code := r.URL.Query().Get("code")
	if code == "" {
		respErr(w, fmt.Errorf("code query param is not provided"))
		return
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	url := fmt.Sprintf("https://oauth.vk.com/access_token?grant_type=authorization_code&code=%s&redirect_uri=%s&client_id=%s&client_secret=%s", code, redirectURI, clientID, clientSecret)
	req, _ := http.NewRequest("POST", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		respErr(w, err)
		return
	}
	defer resp.Body.Close()
	token := struct {
		AccessToken string `json:"access_token"`
	}{}
	bytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bytes, &token)
	url = fmt.Sprintf("https://api.vk.com/method/%s?v=5.124&access_token=%s", "users.get", token.AccessToken)
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		respErr(w, err)
		return
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		respErr(w, err)
		return
	}
	defer resp.Body.Close()
	bytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		respErr(w, err)
		return
	}
	fmt.Fprint(w, string(bytes))
}

func respErr(w http.ResponseWriter, err error) {
	_, er := io.WriteString(w, err.Error())
	if er != nil {
		log.Println(err)
	}
}
