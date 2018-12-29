package maps

import (
	"log"
	"net/http"
	"time"

	"github.com/TukTuk/core"

	gmaps "googlemaps.github.io/maps"
)

var MapsClient *GMClient

type GMClient struct {
	Cfg        *core.Config
	Client     *gmaps.Client
	HttpClient *http.Client
}

func InitMaps(cfg *core.Config) {
	client, err := gmaps.NewClient(gmaps.WithAPIKey(cfg.Maps.ApiKey))
	if err != nil {
		log.Fatal("[InitMaps][Error] Err in creating new gmaps client", err)
	}

	MapsClient = &GMClient{
		Cfg:    cfg,
		Client: client,
		HttpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	log.Println("maps initialized")
}
