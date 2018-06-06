package main

import (
	"net/http"
	"time"

	"github.com/ContaAzul/api-examples/go/auth"
	"github.com/ContaAzul/api-examples/go/product"
	"github.com/apex/log"
	"github.com/gorilla/mux"
)

var (
	// ClientID represents your client id
	ClientID = "PUT YOUR CLIENT_ID HERE"

	// ClientSecret represents your client secret
	ClientSecret = "PUT YOUR CLIENT_SECRET HERE"

	// RedirectURI represents your redirect_uri
	RedirectURI = "PUT YOUR REDIRECT URI HERE"
)

func main() {
	log.Info("initializing")

	mux := mux.NewRouter()

	handler := mux

	mux.Methods(http.MethodGet).Path("/login").HandlerFunc(auth.Authorize(ClientID, RedirectURI))

	mux.Methods(http.MethodGet).Path("/callback").HandlerFunc(auth.Callback(ClientID, ClientSecret, RedirectURI))

	mux.Methods(http.MethodGet).Path("/refresh_token").HandlerFunc(auth.Refresh(ClientID, ClientSecret))

	mux.Methods(http.MethodGet).Path("/list_products").HandlerFunc(product.List())

	mux.Methods(http.MethodGet).Path("/delete_product").HandlerFunc(product.Delete())

	srv := &http.Server{
		Addr:         ":8888",
		Handler:      handler,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
	}

	log.Info("starting")
	if err := srv.ListenAndServe(); err != nil {
		log.WithError(err).Fatal("server fail")
	}
}
