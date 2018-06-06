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

Then access the example in the browser in `http://localhost:/8888`

## Dependencies

This example uses some dependencies to assist some steps :

  - [gorilla/mux](https://github.com/gorilla/mux) : provide routing for go
