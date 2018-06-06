package auth

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/apex/log"
)

// Authorize the application
func Authorize(clientID string, redirectURI string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//must generate a state string and store in cookies or something
		var state = "some-generated-state"

		// your application requests authorization
		var scope = "sales"
		var URL = "https://api.contaazul.com/auth/authorize?client_id=%s&scope=%s&redirect_uri=%s&state=%s"

		http.Redirect(w, r, fmt.Sprintf(URL, clientID, scope, redirectURI, state), 200)
	}
}

// Callback of the authorization
func Callback(clientID string, clientSecret string, redirectURI string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// must get state from stored cookies or something
		var storedState = "some-generated-state"

		defer r.Body.Close() // nolint: errcheck
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.WithError(err).Error("cannot read request parameters")
			http.Redirect(w, r, fmt.Sprintf("/#?error=%s", "'cannot read request parameters'"), 400)
			return
		}

		var parameter struct {
			Code  string `json:"code"`
			State string `json:"state"`
		}
		err = json.Unmarshal(bodyBytes, &parameter)
		if err != nil || parameter.State != storedState {
			log.WithError(err).Error("cannot unmarshal request parameters")
			http.Redirect(w, r, fmt.Sprintf("/#?error=%s", "'cannot unmarshal request parameters'"), 400)
			return
		}

		var form struct {
			Code        string `json:"code"`
			RedirectURI string `json:"redirect_uri"`
			GrantType   string `json:"grant_type"`
		}
		form.Code = parameter.Code
		form.RedirectURI = redirectURI
		form.GrantType = "authorization_code"
		redirectForm, err := json.Marshal(form)
		if nil != err {
			log.WithError(err).Error("cannot marshal form")
			http.Redirect(w, r, fmt.Sprintf("/#?error=%s", "'cannot marshal form'"), 400)
			return
		}

		req, err := http.NewRequest(http.MethodPost, "https://api.contaazul.com/oauth2/token", bytes.NewBuffer(redirectForm))
		if err != nil {
			log.WithError(err).Error("cannot create request")
			http.Redirect(w, r, fmt.Sprintf("/#?error=%s", "'cannot create request'"), 400)
			return
		}

		var data = clientID + clientSecret
		dataEncoded := b64.StdEncoding.EncodeToString([]byte(data))
		req.Header.Set("Authorization", "Basic "+dataEncoded)

		client := &http.Client{Timeout: time.Second * 2}
		resp, err := client.Do(req)
		if err != nil {
			log.WithError(err).Error("cannot post to authorize")
			http.Redirect(w, r, fmt.Sprintf("/#?error=%s", "'cannot post to authorize'"), 400)
			return
		}

		if resp.StatusCode != 200 {
			http.Redirect(w, r, fmt.Sprintf("/#?error=%s", "'invalid_token'"), 400)
			return
		}

		defer resp.Body.Close() // nolint: errcheck
		bodyBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.WithError(err).Error("cannot read response parameters")
			http.Redirect(w, r, fmt.Sprintf("/#?error=%s", "'cannot read response parameters'"), 400)
			return
		}

		var response struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		}
		err = json.Unmarshal(bodyBytes, &response)
		if err != nil {
			log.WithError(err).Error("cannot unmarshal response parameters")
			http.Redirect(w, r, fmt.Sprintf("/#?error=%s", "'cannot unmarshal response parameters'"), 400)
			return
		}
		// we can also pass the token to the browser to make requests from there
		http.Redirect(w, r, fmt.Sprintf("#?access_token=%s&refresh_token=%s", response.AccessToken, response.RefreshToken), 200)
	}
}

// Refresh the token
func Refresh(clientID string, clientSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close() // nolint: errcheck
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.WithError(err).Error("cannot read request parameters")
			http.Redirect(w, r, fmt.Sprintf("/#?error=%s", "'cannot read request parameters'"), 400)
			return
		}

		var parameter struct {
			RefreshToken string `json:"refresh_token"`
			GrantType    string `json:"grant_type"`
		}
		err = json.Unmarshal(bodyBytes, &parameter)
		if err != nil {
			log.WithError(err).Error("cannot unmarshal request parameters")
			http.Redirect(w, r, fmt.Sprintf("/#?error=%s", "'cannot unmarshal request parameters'"), 400)
			return
		}

		parameter.GrantType = "refresh_token"
		redirectForm, err := json.Marshal(parameter)
		if nil != err {
			log.WithError(err).Error("cannot marshal form")
			http.Redirect(w, r, fmt.Sprintf("/#?error=%s", "'cannot marshal form'"), 400)
			return
		}

		req, err := http.NewRequest(http.MethodPost, "https://api.contaazul.com/oauth2/token", bytes.NewBuffer(redirectForm))
		if err != nil {
			log.WithError(err).Error("cannot create request")
			http.Redirect(w, r, fmt.Sprintf("/#?error=%s", "'cannot create request'"), 400)
			return
		}

		var data = clientID + clientSecret
		dataEncoded := b64.StdEncoding.EncodeToString([]byte(data))
		req.Header.Set("Authorization", "Basic "+dataEncoded)

		client := &http.Client{Timeout: time.Second * 2}
		resp, err := client.Do(req)
		if err != nil {
			log.WithError(err).Error("cannot post to refresh token")
			http.Redirect(w, r, fmt.Sprintf("/#?error=%s", "'cannot post to refresh token'"), 400)
			return
		}

		if resp.StatusCode != 200 {
			log.Error("cannot refresh token")
			http.Redirect(w, r, fmt.Sprintf("/#?error=%s", "'cannot refresh token'"), 400)
			return
		}

		defer resp.Body.Close() // nolint: errcheck
		bodyBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.WithError(err).Error("cannot read response parameters")
			http.Redirect(w, r, fmt.Sprintf("/#?error=%s", "'cannot read response parameters'"), 400)
			return
		}

		var response struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		}
		err = json.Unmarshal(bodyBytes, &response)
		if err != nil {
			log.WithError(err).Error("cannot unmarshal response parameters")
			http.Redirect(w, r, fmt.Sprintf("/#?error=%s", "'cannot unmarshal response parameters'"), 400)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("#?access_token=%s&refresh_token=%s", response.AccessToken, response.RefreshToken), 200)

	}
}
