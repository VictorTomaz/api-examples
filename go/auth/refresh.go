package auth

// Refresh the token
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

// Refresh the token authorization
func Refresh(clientID string, clientSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close() // nolint: errcheck
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.WithError(err).Error("cannot read request parameters")
			http.Redirect(w, r, fmt.Sprintf("/#?error=%s", "'cannot read request parameters'"), 400)
			return
		}

		var parameters struct {
			RefreshToken string `json:"refresh_token"`
			GrantType    string `json:"grant_type"`
		}
		err = json.Unmarshal(bodyBytes, &parameters)
		if err != nil {
			log.WithError(err).Error("cannot unmarshal request parameters")
			http.Redirect(w, r, fmt.Sprintf("/#?error=%s", "'cannot unmarshal request parameters'"), 400)
			return
		}

		parameters.GrantType = "refresh_token"
		redirectForm, err := json.Marshal(parameters)
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
