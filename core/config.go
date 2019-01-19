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
	Payment  PaymentConfig
	Ride     RideConfig
}

type RideConfig struct {
	DriverArrival     float64
	DriverArrived     float64
	RideRequestRadius float64
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

type PaymentConfig struct {
	Hostname string
	Timeout  time.Duration
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
		Payment: PaymentConfig{
			Hostname: PAYMENT_STAGING,
			Timeout:  PAYMENT_TIMEOUT,
		},
		Ride: RideConfig{
			DriverArrival: DRIVER_ARRIVAL,
			DriverArrived: DRIVER_ARRIVED,
		},
	}
	return CF
}
