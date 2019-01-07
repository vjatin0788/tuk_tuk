package firebase

import (
	"log"

	"github.com/TukTuk/core"
	fcm "github.com/appleboy/go-fcm"
)

var FClient *FireBase

type FireBase struct {
	FBaseClient *fcm.Client
	Timeout     *uint
}

func InitFireBase(cfg *core.Config) {

	client, err := fcm.NewClient(cfg.FireBase.ApiKey)
	if err != nil {
		log.Fatal("Error in init firebase")
	}

	FClient = &FireBase{
		FBaseClient: client,
		Timeout:     &cfg.FireBase.Timeout,
	}
}
