package main

import (
	"log"
	"net/http"

	"github.com/TukTuk/errs"
	"github.com/TukTuk/firebase"

	"github.com/TukTuk/core"
	"github.com/TukTuk/fulfilment"

	"github.com/TukTuk/authentication"
	"github.com/TukTuk/maps"
	"github.com/TukTuk/model"

	"github.com/TukTuk/api"
)

func main() {

	cfg := core.InitConfig()

	api.InitApiMod()
	api.Api.InitHandler()

	err := model.InitDatabase()
	if err != nil {
		log.Fatal("DB Initialization failed")
	}

	firebase.InitFireBase(cfg)

	authentication.InitAuth()
	fulfilment.InitFF(cfg)

	maps.InitMaps(cfg)
	errs.InitError()

	log.Printf("serving at localhost:%d", 8000)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
