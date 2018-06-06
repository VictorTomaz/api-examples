# ContaAzul authorization example - Golang

## Getting started

In order to run this example, you will need to install Go. On the [official Go website](https://golang.org/) you can learn how to install it.
You'll need to install the [Dep](https://github.com/golang/dep) as well.

With Go and Dep installed, clone this repository, navigate to the _go_ directory and run :

```
  $ go get
  $ dep ensure
```

This command will download all the dependencies required by this example.

## Configuring

Before running the example you must change your credentials in the `main.go` file.
The variables that must be filled with your information are :

 - client_id
 - client_secret
 - redirect_uri

To get this information you'll need to have an API Credential.
For more information about how to get an API Credential go to [http://developers.contaazul.com](http://developers.contaazul.com)

## Running the example

To run the example just execute this command in the _go_ directory :

`$ go run main.go`

The endpoints available are :

- **GET** /login - redirects to the ContaAzul authorization page
- **GET** /callback - the callback of the login
- **GET** /refresh_token - refresh the token
- **GET** /list_products - list all products with pagination
- **GET** /delete_product - delete a product by id

This example doesn't have frontend example yet.

## Dependencies

This example uses some dependencies to assist some steps :

  - [gorilla/mux](https://github.com/gorilla/mux) : provide routing for go
  - [apex/log](https://github.com/apex/log) : provide a better logging for go
