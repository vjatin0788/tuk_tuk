package main

import (
	"log"
	"net/http"

	"github.com/TukTuk/authentication"
	"github.com/TukTuk/model"

	"github.com/TukTuk/api"
)

func main() {
	api.InitApiMod()
	api.Api.InitHandler()

	err := model.InitDatabase()
	if err != nil {
		log.Fatal("DB Initialization failed")
	}

	authentication.InitAuth()

	log.Printf("serving at localhost:%d", 8000)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
