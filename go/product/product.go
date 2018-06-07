package product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/apex/log"
)

// List the products
func List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close() // nolint: errcheck
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.WithError(err).Error("cannot read request parameters")
			http.Redirect(w, r, fmt.Sprintf("/#?error=%s", "'cannot read request parameters'"), 400)
			return
		}

		var parameters struct {
			AccessToken string `json:"access_token"`
			Page        string `json:"page"`
			Size        string `json:"size"`
		}
		err = json.Unmarshal(bodyBytes, &parameters)
		if err != nil {
			log.WithError(err).Error("cannot unmarshal request parameters")
			http.Redirect(w, r, fmt.Sprintf("/#?error=%s", "'cannot unmarshal request parameters'"), 400)
			return
		}

		var APIURL = "https://api.contaazul.com/v1/products?page=%s&size=%s"
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(APIURL, parameters.Page, parameters.Size), nil)
		if err != nil {
			log.WithError(err).Error("cannot create request")
			http.Redirect(w, r, fmt.Sprintf("/#?error=%s", "'cannot create request'"), 400)
			return
		}
		req.Header.Set("Authorization", "Bearer "+parameters.AccessToken)

		client := &http.Client{Timeout: time.Second * 2}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != 200 {
			log.WithError(err).Error("cannot list products")
			http.Redirect(w, r, fmt.Sprintf("/#?error=%s", "'list products'"), 400)
			return
		}
		defer resp.Body.Close() // nolint: errcheck
		bodyBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.WithError(err).Error("cannot read product list bytes")
			http.Redirect(w, r, fmt.Sprintf("/#?error=%s", "'cannot read product list bytes'"), 400)
			return
		}
		w.Write(bodyBytes)
	}
}

// Delete a product
func Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close() // nolint: errcheck
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.WithError(err).Error("cannot read request parameters")
			http.Redirect(w, r, fmt.Sprintf("/#?error=%s", "'cannot read request parameters'"), 400)
			return
		}

		var parameters struct {
			AccessToken string `json:"access_token"`
			ID          string `json:"id"`
		}
		err = json.Unmarshal(bodyBytes, &parameters)
		if err != nil {
			log.WithError(err).Error("cannot unmarshal request parameters")
			http.Redirect(w, r, fmt.Sprintf("/#?error=%s", "'cannot unmarshal request parameters'"), 400)
			return
		}

		var APIURL = "https://api.contaazul.com/v1/products/%s"
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf(APIURL, parameters.ID), nil)
		if err != nil {
			log.WithError(err).Error("cannot create request")
			http.Redirect(w, r, fmt.Sprintf("/#?error=%s", "'cannot create request'"), 400)
			return
		}
		req.Header.Set("Authorization", "Bearer "+parameters.AccessToken)

		client := &http.Client{Timeout: time.Second * 2}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != 204 {
			log.WithError(err).Error("cannot delete product")
			http.Redirect(w, r, fmt.Sprintf("/#?error=%s", "'cannot delete product'"), 400)
			return
		}
		w.Write([]byte("ok"))
	}
}
