package firebase

import (
	fcm "github.com/NaySoftware/go-fcm"
	"github.com/TukTuk/core"
)

var FClient *FireBase

type FireBase struct {
	FBaseClient *fcm.FcmClient
	Ids         []string
	Xds         []string
}

func InitFireBase(cfg *core.Config) {

	client := fcm.NewFcmClient(cfg.FireBase.ApiKey)

	FClient = &FireBase{
		FBaseClient: client,
	}
}
