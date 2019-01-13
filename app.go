package main

import (
	"log"
	"net/http"
	"os"

	"github.com/TukTuk/payment"

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

	f, err := os.OpenFile("access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	log.Println("Log init")

	cfg := core.InitConfig()

	api.InitApiMod()
	api.Api.InitHandler()

	err = model.InitDatabase()
	if err != nil {
		log.Fatal("DB Initialization failed")
	}

	firebase.InitFireBase(cfg)

	authentication.InitAuth()
	fulfilment.InitFF(cfg)

	maps.InitMaps(cfg)
	errs.InitError()

	payment.InitPayment(cfg)

	log.Printf("serving at localhost:%d", 8000)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
