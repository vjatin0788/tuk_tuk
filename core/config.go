package core

import (
	"time"

	"github.com/TukTuk/common"
)

var CF *Config

type Config struct {
	Server   ServerConfig
	Maps     GMaps
	FireBase FBase
}

type GMaps struct {
	Hostname string
	ApiKey   string
}

type FBase struct {
	ApiKey  string
	Timeout uint
}

type ServerConfig struct {
	RideRequestTime time.Duration
}

func InitConfig() *Config {
	CF = &Config{
		Maps: GMaps{
			Hostname: GMAPS_SERVICE_HOSTNAME,
			ApiKey:   common.API_KEY,
		},
		FireBase: FBase{
			ApiKey:  FIREBASE_KEY,
			Timeout: FIREBASE_TIMEOUT,
		},
		Server: ServerConfig{
			RideRequestTime: RIDE_REQUEST_TIME,
		},
	}
	return CF
}
